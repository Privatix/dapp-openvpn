package openvpn

import (
	"fmt"

	"github.com/privatix/dappctrl/util/log"
)

// Flow is a slice of Operators that can be applied in sequence.
type Flow []*Operator

// Run executes the Operators flow.
func (flow Flow) Run(in *OpenVPN, logger log.Logger) error {
	rollback := func(f Flow) {
		for _, v := range f {
			defer v.Cancel(in)
		}
	}

	var err error
	for i, m := range flow {
		err = m.Run(in)

		if err != nil {
			logger.Warn(fmt.Sprintf("failed to execute '%v' operation",
				m.Name))
			rollback(flow[:i+1])
			break
		}
		logger.Info(fmt.Sprintf("'%v' operation was successfully executed",
			m.Name))
	}
	return err
}
