package connstr

import (
	"flag"
)

const (
	name = "connstr"
)

// ConnFlag is a object that is used to process the connect flag.
type ConnFlag struct {
	run func()
	val *string
}

// NewConnFlag initializes the object to process the connect flag.
func NewConnFlag() *ConnFlag {
	return &ConnFlag{
		run: func() {},
		val: flag.String("connstr",
			"user=postgres dbname=dappctrl sslmode=disable",
			"PostgreSQL connection string"),
	}
}

// Name returns flag name.
func (c *ConnFlag) Name() string {
	return name
}

// Value returns value of flag.
func (c *ConnFlag) Value() interface{} {
	return c.val
}

// Process performs processing of the flag.
func (c *ConnFlag) Process() error {
	return nil
}
