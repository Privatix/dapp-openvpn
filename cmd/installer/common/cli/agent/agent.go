package agent

import (
	"flag"
)

const (
	name = "agent"
)

// WhetherAgentFlag is a object that is used to process the agent flag.
type WhetherAgentFlag struct {
	run func()
	val *bool
}

// NewWhetherAgentFlag initializes the object to process the agent flag.
func NewWhetherAgentFlag() *WhetherAgentFlag {
	return &WhetherAgentFlag{
		run: func() {},
		val: flag.Bool("agent", false, "Whether to install agent"),
	}
}

// Name returns flag name.
func (c *WhetherAgentFlag) Name() string {
	return name
}

// Value returns value of flag.
func (c *WhetherAgentFlag) Value() interface{} {
	return c.val
}

// Process performs processing of the flag.
func (c *WhetherAgentFlag) Process() error {
	return nil
}
