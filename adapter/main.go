package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/privatix/dappctrl/data"
	"github.com/privatix/dappctrl/nat"
	"github.com/privatix/dappctrl/sess"
	"github.com/privatix/dappctrl/util"
	"github.com/privatix/dappctrl/util/log"
	"github.com/privatix/dappctrl/version"

	"github.com/privatix/dapp-openvpn/adapter/config"
	vpndata "github.com/privatix/dapp-openvpn/adapter/data"
	"github.com/privatix/dapp-openvpn/adapter/mon"
	"github.com/privatix/dapp-openvpn/adapter/msg"
	"github.com/privatix/dapp-openvpn/adapter/prepare"
	"github.com/privatix/dapp-openvpn/adapter/tc"
)

// Values for versioning.
var (
	Commit  string
	Version string
)

var (
	conf    *config.Config
	channel string
	logger  log.Logger
	tctrl   *tc.TrafficControl
	fatal   = make(chan string)
	sesscl  *sess.Client
)

func createLogger() (log.Logger, io.Closer, error) {
	elog, err := log.NewStderrLogger(conf.FileLog.WriterConfig)
	if err != nil {
		return nil, nil, err
	}

	flog, closer, err := log.NewFileLogger(conf.FileLog)
	if err != nil {
		return nil, nil, err
	}

	logger := log.NewMultiLogger(elog, flog)

	return logger, closer, nil
}

func main() {
	v := flag.Bool("version", false, "Prints current dappctrl version")

	fconfig := flag.String(
		"config", "adapter.config.json", "Configuration file")
	flag.Parse()

	version.Print(*v, Commit, Version)

	conf = config.NewConfig()
	if err := util.ReadJSONFile(*fconfig, &conf); err != nil {
		panic(fmt.Sprintf("failed to read configuration: %s\n", err))
	}

	var err error

	var closer io.Closer
	logger, closer, err = createLogger()
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %s", err))
	}
	defer closer.Close()

	sesscl, err = sess.Dial(context.Background(), conf.Sess.Endpoint,
		conf.Sess.Origin, conf.Sess.Product, conf.Sess.Password)
	if err != nil {
		panic("failed to connect to session server: " + err.Error())
	}

	tctrl = tc.NewTrafficControl(conf.TC, logger)

	switch os.Getenv("script_type") {
	case "user-pass-verify":
		handleAuth()
	case "client-connect":
		handleConnect()
	case "client-disconnect":
		handleDisconnect()
	default:
		handleMonitor(*fconfig)
	}
}

func handleAuth() {
	logger := logger.Add("method", "handleAuth")
	user, pass := getCreds()

	if err := sesscl.AuthClient(user, pass); err != nil {
		logger.Fatal("failed to auth: " + err.Error())
	}

	if cn := commonNameOrEmpty(); len(cn) != 0 {
		storeChannel(cn, user)
	}
	storeChannel(user, user) // Needed when using username-as-common-name.
}

func handleConnect() {
	logger := logger.Add("method", "handleConnect")

	port, err := strconv.Atoi(os.Getenv("trusted_port"))
	if err != nil || port <= 0 || port > 0xFFFF {
		logger.Fatal("bad trusted_port value")
	}

	var offer *data.Offering
	offer, err = sesscl.StartSession(os.Getenv("trusted_ip"), loadChannel(), uint16(port))
	if err != nil {
		logger.Fatal("failed to start session: " + err.Error())
	}

	if len(channel) != 0 || offer.AdditionalParams == nil {
		return
	}

	var params vpndata.OfferingParams
	err = json.Unmarshal(offer.AdditionalParams, &params)
	if err != nil {
		logger.Add("offering_params", offer.AdditionalParams).Fatal(
			"failed to unmarshal offering params: " + err.Error())
	}

	err = tctrl.SetRateLimit(os.Getenv("dev"),
		os.Getenv("ifconfig_pool_remote_ip"),
		params.MinUploadMbits, params.MinDownloadMbits)
	if err != nil {
		logger.Fatal("failed to set rate limit: " + err.Error())
	}
}

func handleDisconnect() {
	logger := logger.Add("method", "handleDisconnect")

	down, err := strconv.ParseUint(os.Getenv("bytes_sent"), 10, 64)
	if err != nil || down < 0 {
		panic("bad bytes_sent value")
	}

	up, err := strconv.ParseUint(os.Getenv("bytes_received"), 10, 64)
	if err != nil || up < 0 {
		panic("bad bytes_received value")
	}

	err = sesscl.StopSession(loadChannel())
	if err != nil {
		logger.Fatal("failed to stop session: " + err.Error())
	}

	err = tctrl.UnsetRateLimit(os.Getenv("dev"),
		os.Getenv("ifconfig_pool_remote_ip"))
	if err != nil {
		logger.Fatal("failed to unset rate limit: " + err.Error())
	}
}

func handleMonitor(confFile string) {
	if conf.ClientMode {
		handleClientMonitor()
	} else {
		handleAgentMonitor(confFile)
	}
}

type sessionHandler struct{}

func (h sessionHandler) StartSession(ch string) bool {
	logger := logger.Add("method", "handleMonStarted", "channel", ch)

	if _, err := sesscl.StartSession(os.Getenv("trusted_ip"), ch, 0); err != nil {
		logger.Fatal("failed to start session: " + err.Error())
		return false
	}

	return true
}

func (h sessionHandler) UpdateSession(ch string, up, down uint64) bool {
	logger := logger.Add("method", "handleMonByteCount",
		"channel", ch, "up", up, "down", down)

	err := sesscl.UpdateSession(ch, down+up)
	if err != nil {
		logger.Fatal("failed to update session: " + err.Error())
		return false
	}

	return true
}

func (sessionHandler) StopSession(ch string) bool {
	logger := logger.Add("method", "handleMonStopped", "channel", ch)

	err := sesscl.StopSession(ch)
	if err != nil {
		logger.Fatal("failed to stop session: " + err.Error())
		return false
	}

	return true
}

func openExtPort(ctx context.Context, network string, port int) {
	logger := logger.Add("method", "openExtPort")

	service, err := nat.Parse(conf.NAT.Config)
	if err != nil {
		logger.Debug("failed to create nat interface, error: " +
			err.Error())
		return
	}

	service.DeleteMapping(network, port, port)

	if err := nat.Map(ctx, conf.NAT.Config, logger, service, network,
		port, port, "openvpn"); err != nil {
		logger.Debug("failed to mapping ports to nat interface," +
			" error: " + err.Error())
	}
}

func handlePusher(ctx context.Context, dir string) {
	logger := logger.Add("method", "handlePusher", "directory", dir)

	pusher := msg.NewPusher(conf.Pusher, logger,
		func(config map[string]string) error {
			return sesscl.SetProductConfig(config)
		})

	params, err := pusher.VpnParams()
	if err != nil {
		return
	}

	network := "tcp"
	netPort := 443

	if proto, ok := params["proto"]; ok {
		if len(proto) > 3 {
			proto = proto[:3]
		}
		network = proto
	}

	if portStr, ok := params["port"]; ok {
		port, err := strconv.Atoi(portStr)
		if err == nil {
			netPort = port
		}
	}

	go openExtPort(ctx, network, netPort)

	if msg.IsDone(dir) {
		return
	}

	err = pusher.PushConfiguration(ctx, params)
	if err != nil {
		logger.Error("failed to push app config to" +
			" dappctrl")
		return
	}

	err = msg.Done(dir)
	if err == nil {
		return
	}

	logger.Add("file", msg.PushedFile).Error(
		"failed to save file in directory: " + err.Error())
}

func handleAgentMonitor(confFile string) {
	dir := filepath.Dir(confFile)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handlePusher(ctx, dir)

	monitor := mon.NewMonitor(conf.Monitor, logger, &sessionHandler{}, channel)
	go func() {
		fatal <- fmt.Sprintf("failed to monitor vpn traffic: %s",
			monitor.MonitorTraffic(context.Background()))
	}()

	onConnStart := func(channel string) {
		if err := sesscl.ServiceReady(channel); err != nil {
			fatal <- "could not signal that service is ready"
		}
	}
	onConnStop := func(channel string) {
		// Call stop session to signal ctrl that we all done to terminate a service.
		if !(sessionHandler{}).StopSession(channel) {
			fatal <- "failed to stop session"
		}
	}

	subscribeAndStartHandlingConnChanges(onConnStart, onConnStop)

	logger.Fatal(<-fatal)
}

var (
	mtx     sync.Mutex
	ovpnCmd *exec.Cmd
)

func handleClientMonitor() {
	if channel := loadActiveChannel(); len(channel) != 0 {
		logger.Warn("interrupted connection detected: " + channel)
		sessionHandler{}.StopSession(channel)
		removeActiveChannel()
	}

	getEndpoint := func(clientKey string) (*data.Endpoint, error) {
		ept, err := sesscl.GetEndpoint(clientKey)
		return ept, err
	}

	stopOvpnAndMonitor := func() {}
	defer stopOvpnAndMonitor()

	onConnStart := func(channel string) {
		if ovpnCmd != nil {
			logger.Warn("requested to start while " +
				"OpenVPN is still running")
			return
		}

		err := prepare.ClientConfig(logger, channel, conf, getEndpoint)
		if err != nil {
			message := "failed to prepare client config: "
			logger.Fatal(message + err.Error())
		}

		ctx, cancel := context.WithCancel(context.Background())
		stopOvpnAndMonitor = cancel
		ovpnCmd = launchOpenVPN(ctx, channel)
	}

	onConnStop := func(channel string) {
		if ovpnCmd == nil {
			logger.Warn("requested to stop while OpenVPN is not running")
			sessionHandler{}.StopSession(channel)
			return
		}

		stopOvpnAndMonitor()
	}

	subscribeAndStartHandlingConnChanges(onConnStart, onConnStop)
}

func subscribeAndStartHandlingConnChanges(onStart, onStop func(string)) {
	ch := make(chan *sess.ConnChangeResult)
	subcl, err := sesscl.ConnChange(ch)
	if err != nil {
		logger.Fatal("failed to subscribe to connection changes")
	}

	go func() {
		select {
		case err := <-subcl.Err():
			logger.Add("err", err).Fatal(
				"unexpected end of subscription to connection changes")
		}
	}()

	for res := range ch {
		logger.Info(fmt.Sprintf("connection change: %v", res))

		mtx.Lock()
		switch res.Status {
		case sess.ConnStart:
			onStart(res.Channel)
		case sess.ConnStop:
			onStop(res.Channel)
		}
		mtx.Unlock()
	}

	logger.Fatal("unexpected end of subscription to connection changes")
}

func launchOpenVPN(ctx context.Context, channel string) *exec.Cmd {
	if len(conf.OpenVPN.Name) == 0 {
		logger.Fatal("no OpenVPN command provided")
	}

	args := append(conf.OpenVPN.Args, "--config")
	args = append(args, filepath.Join(
		conf.OpenVPN.ConfigRoot, channel, "client.ovpn"))

	cmd := exec.Command(conf.OpenVPN.Name, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Fatal("failed to access OpenVPN stdout: " + err.Error())
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Fatal("failed to access OpenVPN stderr: " + err.Error())
	}

	if err := cmd.Start(); err != nil {
		logger.Fatal("failed to launch OpenVPN: " + err.Error())
	}

	storeActiveChannel(channel)

	go func() {
		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		for scanner.Scan() {
			line := "openvpn: " + scanner.Text() + "\n"
			io.WriteString(os.Stderr, line)
		}
		if err := scanner.Err(); err != nil {
			message := "failed to read from openVPN stdout/stderr: "
			logger.Warn(message + err.Error())
		}
		stdout.Close()
		stderr.Close()
	}()

	go func() {
		logger.Warn(fmt.Sprintf("OpenVPN exited: %v", cmd.Wait()))
		sessionHandler{}.StopSession(channel)
		removeActiveChannel()
		mtx.Lock()
		ovpnCmd = nil
		mtx.Unlock()
	}()

	time.Sleep(conf.OpenVPN.StartDelay * time.Millisecond)

	monitor := mon.NewMonitor(conf.Monitor, logger, &sessionHandler{}, channel)
	go func() {
		err := monitor.MonitorTraffic(ctx)
		logger.Warn("failed to monitor vpn traffic: " + err.Error())
	}()

	return cmd
}
