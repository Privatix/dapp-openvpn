package command

import "github.com/privatix/dapp-openvpn/inst/openvpn"

// operator implement the Runner interface in pipeline package.
type operator struct {
	name   string
	run    func(*openvpn.OpenVPN) error
	cancel func(*openvpn.OpenVPN) error
}

// Name returns the operators name.
func (o operator) Name() string {
	return o.name
}

// Run executes the operators run function.
func (o operator) Run(in interface{}) error {
	return o.run(in.(*openvpn.OpenVPN))
}

// Cancel executes the operators cancel function.
func (o operator) Cancel(in interface{}) error {
	return o.cancel(in.(*openvpn.OpenVPN))
}

func newOperator(name string, run func(*openvpn.OpenVPN) error,
	cancel func(*openvpn.OpenVPN) error) operator {
	if cancel == nil {
		cancel = func(in *openvpn.OpenVPN) error { return nil }
	}
	return operator{name: name, run: run, cancel: cancel}
}
