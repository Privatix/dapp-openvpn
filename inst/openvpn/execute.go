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
	Process *os.Process
}

// Start is a start method of executable service.
func (e *execute) Start() {
	go func() error {
		ovpn := filepath.Join(e.Path, path.OpenVPN)
		config := filepath.Join(e.Path, path.RoleConfig(e.Role))
		cmd := exec.Command(ovpn, "--config", config)

		if err := cmd.Start(); err != nil {
			return err
		}

		e.Process = cmd.Process

		return cmd.Wait()
	}()
}

// Run is a run method of executable service.
func (e *execute) Run() {
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	defer close(interrupt)

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
