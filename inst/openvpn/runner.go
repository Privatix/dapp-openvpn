package openvpn

import (
	"github.com/privatix/dappctrl/util/log"
)

// Operator aliases a runner function.
type Operator struct {
	Run    func(*OpenVPN) error
	Cancel func(*OpenVPN) error
}

// Flow is a slice of Operators that can be applied in sequence.
type Flow []*Operator

// Exec executes the operators run function.
func (o Operator) Exec(in *OpenVPN) error {
	return o.Run(in)
}

// Run executes the flow operators runner function.
func (flow Flow) Run(in *OpenVPN, logger log.Logger) error {
	rollback := func(f Flow) {
		for _, v := range f {
			defer v.Cancel(in)
		}
	}

	var err error
	for i, m := range flow {
		err = m.Exec(in)

		if err != nil {
			logger.Error(err.Error())
			rollback(flow[:i+1])
			break
		}
	}
	return err
}

// NewOperator creates a new operator instance.
func NewOperator(run func(*OpenVPN) error, cancel func(*OpenVPN) error) *Operator {
	if cancel == nil {
		cancel = func(in *OpenVPN) error { return nil }
	}
	return &Operator{Run: run, Cancel: cancel}
}
