package mon

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/privatix/dappctrl/util/log"
)

// Config is a configuration for OpenVPN monitor.
type Config struct {
	Addr            string
	ByteCountPeriod uint // In seconds.
	CmdApplyTimeout uint // In seconds.
	CmdRetryTimeout uint // In seconds.
}

// NewConfig creates a default configuration for OpenVPN monitor.
func NewConfig() *Config {
	return &Config{
		Addr:            "localhost:7505",
		ByteCountPeriod: 5,
		CmdApplyTimeout: 10,
		CmdRetryTimeout: 3,
	}
}

type client struct {
	channel    string
	commonName string
}

// Monitor is an OpenVPN monitor for observation of consumed VPN traffic and
// for killing client VPN sessions.
type Monitor struct {
	conf            *Config
	logger          log.Logger
	handleSession   HandleSessionFunc
	channel         string // Client mode channel (empty in server mode).
	conn            net.Conn
	mtx             sync.Mutex // To guard writing.
	clients         map[uint]client
	clientConnected bool
	out             *bufio.Reader // Openvpn output.
}

// Session events.
const (
	SessionStarted   = iota // For client mode only.
	SessionStopped   = iota // For client mode only.
	SessionByteCount = iota
)

// HandleSessionFunc is a session event handler. If it returns false in server
// mode, then the monitor kills the corresponding session.
type HandleSessionFunc func(ch string, event int, up, down uint64) bool

// NewMonitor creates a new OpenVPN monitor.
func NewMonitor(conf *Config, logger log.Logger,
	handleSession HandleSessionFunc, channel string) *Monitor {
	return &Monitor{
		conf:          conf,
		logger:        logger,
		handleSession: handleSession,
		channel:       channel,
	}
}

// Close immediately closes the monitor making MonitorTraffic() to return.
func (m *Monitor) Close() error {
	if m.conn != nil {
		return m.conn.Close()
	}
	return nil
}

// MonitorTraffic connects to OpenVPN management interfaces and starts
// monitoring VPN traffic.
func (m *Monitor) MonitorTraffic(ctx context.Context) error {
	logger := m.logger.Add("method", "MonitorTraffic")
	logger.Info("dapp-openvpn monitor started")

	var err error
	if m.conn, err = net.Dial("tcp", m.conf.Addr); err != nil {
		return err
	}
	defer m.conn.Close()

	m.out = bufio.NewReader(m.conn)

	if err := m.initConn(); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			logger.Debug("context cancelled, exiting")
			return ErrMonitoringCancelled
		default:
			str, err := m.out.ReadString('\n')
			if err != nil {
				return err
			}

			if err = m.processReply(str); err != nil {
				return err
			}
		}
	}
}

func (m *Monitor) writeAndWaitForSuccess(cmd string) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	ret := make(chan error)

	stopch := make(chan struct{})

	go func() {
		period := time.Duration(m.conf.CmdRetryTimeout) * time.Second
		err := m.writeCommandEach(cmd, period, stopch)
		ret <- err
	}()

	go func() {
		timeout := time.Duration(m.conf.CmdApplyTimeout) * time.Second
		// Prefix "SUCCESS: " indicates that openvpn received the command.
		err := m.lookForPrefixOutput(timeout, "SUCCESS: ")
		ret <- err
		stopch <- struct{}{}
	}()

	return <-ret
}

func (m *Monitor) write(cmd string) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return m.writeCommand(cmd)
}

func (m *Monitor) writeCommandEach(cmd string, period time.Duration, exit chan struct{}) error {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := m.writeCommand(cmd)
			if err != nil {
				return err
			}
		case <-exit:
			return nil
		}
	}
}

func (m *Monitor) writeCommand(cmd string) error {
	_, err := m.conn.Write([]byte(cmd + "\n"))
	return err
}

func (m *Monitor) lookForPrefixOutput(timeout time.Duration, prefix string) error {
	logger := m.logger.Add("method", "lookForPrefixOutput",
		"timeout", timeout, "prefix", prefix)
	logger.Debug("looking for prefixe: " + prefix)

	failchan := time.After(timeout)
	for {
		select {
		case <-failchan:
			logger.Error("looking for prefix output timeout")
			return ErrCmdReceiveTimeout
		default:
			out, err := m.out.ReadString('\n')
			logger.Debug("out: " + out)
			if err != nil || strings.HasPrefix(out, prefix) {
				return err
			}
		}
	}
}

func (m *Monitor) requestClients() error {
	m.logger.Info("requesting updated client list")
	return m.write("status 2")
}

func (m *Monitor) setByteCountPeriod() error {
	return m.writeAndWaitForSuccess(fmt.Sprintf("bytecount %d", m.conf.ByteCountPeriod))
}

func (m *Monitor) killSession(cn string) error {
	return m.write(fmt.Sprintf("kill %s", cn))
}

func (m *Monitor) initConn() error {
	if err := m.setByteCountPeriod(); err != nil {
		return err
	}

	if len(m.channel) == 0 {
		if err := m.requestClients(); err != nil {
			return err
		}
	} else {
		if err := m.writeAndWaitForSuccess("state on"); err != nil {
			return err
		}

		if err := m.writeAndWaitForSuccess("hold release"); err != nil {
			return err
		}

		// On Windows OpenVPN doesn't react on first `hold release`
		// for some reason. This needs to be further investigated.
		if runtime.GOOS == "windows" {
			time.Sleep(time.Millisecond * 100)
			if err := m.writeAndWaitForSuccess("hold release"); err != nil {
				return err
			}
		}
	}

	return nil
}

const (
	prefixClientListHeader  = "HEADER,CLIENT_LIST,"
	prefixClientList        = "CLIENT_LIST,"
	prefixByteCount         = ">BYTECOUNT_CLI:"
	prefixByteCountClient   = ">BYTECOUNT:"
	prefixClientEstablished = ">CLIENT:ESTABLISHED,"
	prefixError             = "ERROR: "
	prefixState             = ">STATE:"
)

func (m *Monitor) processReply(s string) error {
	logger := m.logger.Add("method", "processReply", "reply", s)

	logger.Debug("openvpn raw: " + s)

	if strings.HasPrefix(s, prefixClientListHeader) {
		m.clients = make(map[uint]client)
		return nil
	}

	if strings.HasPrefix(s, prefixClientList) {
		return m.processClientList(s[len(prefixClientList):])
	}

	if strings.HasPrefix(s, prefixByteCount) {
		return m.processByteCount(s[len(prefixByteCount):])
	}

	if strings.HasPrefix(s, prefixByteCountClient) {
		return m.processByteCountClient(s[len(prefixByteCountClient):])
	}

	if strings.HasPrefix(s, prefixClientEstablished) {
		return m.requestClients()
	}

	if strings.HasPrefix(s, prefixState) {
		return m.processState(s[len(prefixState):])
	}

	if strings.HasPrefix(s, prefixError) {
		logger.Error("openvpn error: " + s[len(prefixError):])
	}

	return nil
}

func split(s string) []string {
	return strings.Split(strings.TrimRight(s, "\r\n"), ",")
}

func (m *Monitor) processClientList(s string) error {
	logger := m.logger.Add("method", "processClientList")

	sp := split(s)
	if len(sp) < 10 {
		return ErrServerOutdated
	}

	cid, err := strconv.ParseUint(sp[9], 10, 32)
	if err != nil {
		return err
	}

	m.clients[uint(cid)] = client{sp[8], sp[0]}
	logger.Info(fmt.Sprintf("openvpn client found:"+
		" cid %d, chan %s, cn %s", cid, sp[8], sp[0]))

	return nil
}

func (m *Monitor) processByteCount(s string) error {
	logger := m.logger.Add("method", "processByteCount")

	sp := split(s)

	cid, err := strconv.ParseUint(sp[0], 10, 32)
	if err != nil {
		return err
	}

	down, err := strconv.ParseUint(sp[1], 10, 64)
	if err != nil {
		return err
	}

	up, err := strconv.ParseUint(sp[2], 10, 64)
	if err != nil {
		return err
	}

	cl, ok := m.clients[uint(cid)]
	if !ok {
		return m.requestClients()
	}

	logger.Info(fmt.Sprintf("openvpn byte count for chan %s:"+
		" up %d, down %d", cl.channel, up, down))

	go func() {
		if !m.handleSession(cl.channel, SessionByteCount, up, down) {
			m.killSession(cl.commonName)
		}
	}()

	return nil
}

func (m *Monitor) processByteCountClient(s string) error {
	logger := m.logger.Add("method", "processByteCountClient")

	if !m.clientConnected {
		return nil
	}

	sp := split(s)

	down, err := strconv.ParseUint(sp[0], 10, 64)
	if err != nil {
		return err
	}

	up, err := strconv.ParseUint(sp[1], 10, 64)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf(
		"openvpn byte count: up %d, down %d", up, down))

	go func() {
		m.handleSession(m.channel, SessionByteCount, up, down)
	}()

	return nil
}

func (m *Monitor) processState(s string) error {
	logger := m.logger.Add("method", "processState")

	connected := split(s)[1] == "CONNECTED"

	if m.clientConnected && !connected {
		logger.Warn("disconnected from server")
		go func() {
			m.handleSession(m.channel, SessionStopped, 0, 0)
		}()
		m.clientConnected = false
	} else if !m.clientConnected && connected {
		logger.Warn("connected to server")
		go func() {
			m.handleSession(m.channel, SessionStarted, 0, 0)
		}()
		m.clientConnected = true
	}

	return nil
}
