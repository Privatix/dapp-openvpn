package pipeline

import (
	"encoding/json"
	"fmt"

	"github.com/privatix/dappctrl/util/log"
)

// Runner interface has a pipeline's element functions.
type Runner interface {
	Name() string
	Run(interface{}) error
	Cancel(interface{}) error
}

// Flow is a slice of Runners interface elements
// that can be applied in sequence.
type Flow []Runner

// Run executes the flow elements runner function.
func (flow Flow) Run(in interface{}, logger log.Logger) error {
	rollback := func(f Flow) {
		for _, v := range f {
			defer v.Cancel(in)
		}
	}

	var err error
	for i, m := range flow {
		err = m.Run(in)

		if err != nil {
			object, _ := json.Marshal(in)
			l := logger.Add("object", string(object))
			l.Warn(fmt.Sprintf("failed to execute '%v' operation",
				m.Name()))
			rollback(flow[:i+1])
			break
		}
		logger.Info(fmt.Sprintf("'%v' operation was successfully executed",
			m.Name()))
	}
	return err
}
