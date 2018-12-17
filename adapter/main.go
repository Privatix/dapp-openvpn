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

	"github.com/privatix/dappctrl/sesssrv"
	"github.com/privatix/dappctrl/svc/connector"
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

//go:generate go generate ../vendor/github.com/privatix/dappctrl/data/schema.go

// Values for versioning.
var (
	Commit  string
	Version string
)

var (
	conf    *config.Config
	conn    connector.Connector
	channel string
	logger  log.Logger
	tctrl   *tc.TrafficControl
	fatal   = make(chan string)
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

	conn = connector.NewConnector(conf.Connector, logger)

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
	args := &sesssrv.AuthArgs{ClientID: user, Password: pass}

	err := conn.AuthSession(args)
	if err != nil {
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

	args := &sesssrv.StartArgs{
		ClientID:   loadChannel(),
		ClientIP:   os.Getenv("trusted_ip"),
		ClientPort: uint16(port),
	}

	res, err := conn.StartSession(args)
	if err != nil {
		logger.Fatal("failed to start session: " + err.Error())
	}

	if len(channel) != 0 || res.Offering.AdditionalParams == nil {
		return
	}

	var params vpndata.OfferingParams
	err = json.Unmarshal(res.Offering.AdditionalParams, &params)
	if err != nil {
		logger.Add("offering_params", res.Offering.AdditionalParams).Fatal(
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

	args := &sesssrv.StopArgs{
		ClientID: loadChannel(),
		Units:    down + up,
	}

	err = conn.StopSession(args)
	if err != nil {
		logger.Fatal("failed to stop session: " + err.Error())
	}

	err = tctrl.UnsetRateLimit(os.Getenv("dev"),
		os.Getenv("ifconfig_pool_remote_ip"))
	if err != nil {
		logger.Fatal("failed to unset rate limit: " + err.Error())
	}
}

func handleMonStarted(ch string) bool {
	logger := logger.Add("method", "handleMonStarted", "channel", ch)

	args := &sesssrv.StartArgs{
		ClientID: ch,
	}

	_, err := conn.StartSession(args)
	if err != nil {
		logger.Fatal("failed to start session: " + err.Error())
		return false
	}

	return true
}

func handleMonStopped(ch string, up, down uint64) bool {
	logger := logger.Add("method", "handleMonStopped",
		"channel", ch, "up", up, "down", down)

	args := &sesssrv.StopArgs{
		ClientID: ch,
		Units:    down + up,
	}

	err := conn.StopSession(args)
	if err != nil {
		logger.Fatal("failed to stop session: " + err.Error())
		return false
	}

	return true
}

func handleMonByteCount(ch string, up, down uint64) bool {
	logger := logger.Add("method", "handleMonByteCount",
		"channel", ch, "up", up, "down", down)

	args := &sesssrv.UpdateArgs{
		ClientID: ch,
		Units:    down + up,
	}

	err := conn.UpdateSessionUsage(args)
	if err != nil {
		logger.Error("failed to update session: " + err.Error())
		return false
	}

	return true
}

func handleMonitor(confFile string) {
	if conf.ClientMode {
		handleClientMonitor()
	} else {
		handleAgentMonitor(confFile)
	}
}

func handleSession(ch string, event int, up, down uint64) bool {
	switch event {
	case mon.SessionStarted:
		return handleMonStarted(ch)
	case mon.SessionStopped:
		return handleMonStopped(ch, up, down)
	case mon.SessionByteCount:
		return handleMonByteCount(ch, up, down)
	default:
		return false
	}
}

func handlePusher(ctx context.Context, dir string) {
	logger := logger.Add("method", "handlePusher", "directory", dir)

	pusher := msg.NewPusher(conf.Pusher, logger, conn)

	err := pusher.PushConfiguration(ctx)
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

	if !msg.IsDone(dir) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go handlePusher(ctx, dir)
	}

	monitor := mon.NewMonitor(conf.Monitor, logger, handleSession, channel)
	go func() {
		fatal <- fmt.Sprintf("failed to monitor vpn traffic: %s",
			monitor.MonitorTraffic())
	}()

	logger.Fatal(<-fatal)
}

var (
	mtx     sync.Mutex
	ovpnCmd *exec.Cmd
)

func handleClientMonitor() {
	if channel := loadActiveChannel(); len(channel) != 0 {
		logger.Warn("interrupted connection detected: " + channel)
		handleMonStopped(channel, 0, 0)
		removeActiveChannel()
	}

	for {
		time.Sleep(conf.HeartbeatPeriod * time.Millisecond)

		res, err := conn.Heartbeat()
		if err != nil {
			logger.Error("heartbeat request failed: " + err.Error())
			break
		}

		mtx.Lock()

		switch res.Command {
		case sesssrv.HeartbeatStart:
			if ovpnCmd != nil {
				logger.Warn("requested to start while " +
					"OpenVPN is still running")
				break
			}

			err := prepare.ClientConfig(
				logger, res.Channel, conn, conf)
			if err != nil {
				message := "failed to prepare client config: "
				logger.Fatal(message + err.Error())
			}

			ovpnCmd = launchOpenVPN(res.Channel)

		case sesssrv.HeartbeatStop:
			if ovpnCmd == nil {
				logger.Warn("requested to stop while " +
					"OpenVPN is not running")
				handleMonStopped(res.Channel, 0, 0)
				break
			}

			if err := ovpnCmd.Process.Kill(); err != nil {
				message := "failed to kill OpenVPN: "
				logger.Error(message + err.Error())
			}
		}

		mtx.Unlock()
	}
}

func launchOpenVPN(channel string) *exec.Cmd {
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
		handleMonStopped(channel, 0, 0)
		removeActiveChannel()
		mtx.Lock()
		ovpnCmd = nil
		mtx.Unlock()
	}()

	time.Sleep(conf.OpenVPN.StartDelay * time.Millisecond)

	monitor := mon.NewMonitor(conf.Monitor, logger, handleSession, channel)
	go func() {
		err := monitor.MonitorTraffic()
		logger.Warn("failed to monitor vpn traffic: " + err.Error())
	}()

	return cmd
}
