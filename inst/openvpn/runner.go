package openvpn

import (
	"github.com/privatix/dappctrl/util/log"
)

// Operator aliases a runner function.
type Operator struct {
	Run    func(*OpenVPN) error
	Cancel func(*OpenVPN)
}

// Flow is a slice of Operators that can be applied in sequence.
type Flow []*Operator

// Exec executes the operators run function.
func (o Operator) Exec(in *OpenVPN) error {
	return o.Run(in)
}

// Run executes the flow operators runner function.
func (flow Flow) Run(in *OpenVPN, logger log.Logger) error {
	var err error
	for _, m := range flow {
		err = m.Exec(in)
		if err != nil {
			logger.Error(err.Error())
			m.Cancel(in)
			//todo rollback
			break
		}
	}
	return err
}

// NewOperator creates a new operator instance.
func NewOperator(run func(*OpenVPN) error, cancel func(*OpenVPN)) *Operator {
	if cancel == nil {
		cancel = func(in *OpenVPN) {}
	}
	return &Operator{Run: run, Cancel: cancel}
}
