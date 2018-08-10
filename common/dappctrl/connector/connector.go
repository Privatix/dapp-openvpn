// Package connector implements standard methods for communicating with dappctrl.
package connector

// Connector defines the methods for interacting with a dappctrl.
type Connector interface {
	AuthSession(args interface{}) error
	StartSession(args interface{}) error
	StopSession(args interface{}) error
	UpdateSessionUsage(args interface{}) error
	SetupProductConfiguration(args interface{}) error
}
