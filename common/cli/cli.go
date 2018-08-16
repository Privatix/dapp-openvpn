package cli

// Flag defines the methods for processing flags.
type Flag interface {
	Name() string
	Value() interface{}
	Process() error
}
