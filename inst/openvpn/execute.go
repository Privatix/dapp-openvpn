package openvpn

import (
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/privatix/dapp-openvpn/inst/openvpn/path"
)

type execute struct {
	Path    string
	Role    string
	Type    string
	Process *os.Process
}

// Start is a start method of executable service.
func (e *execute) Start() {
	go func() {
		if err := run(e); err != nil {
			os.Exit(2)
		}
		os.Exit(0)
	}()
}

func run(e *execute) error {
	vpn := filepath.Join(e.Path, path.VPN(e.Type))
	config := filepath.Join(e.Path, path.VPNConfig(e.Type, e.Role))
	args := []string{}
	if e.Type == path.Config.OVPN {
		args = append(args, "--cd", e.Path)
	}
	args = append(args, "--config", config)
	cmd := exec.Command(vpn, args...)

	if err := cmd.Start(); err != nil {
		return err
	}

	e.Process = cmd.Process

	return cmd.Wait()
}

// Run is a run method of executable service.
func (e *execute) Run() {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	defer close(interrupt)

	go run(e)

	for {
		select {
		case <-interrupt:
			if e.Process != nil {
				e.Process.Kill()
			}
			break
		}
	}
}

// Stop is a stop method of executable service.
func (e *execute) Stop() {
	if e.Process != nil {
		e.Process.Kill()
	}
}
