// Package logger provides standard logger interface for applications.
package logger

// Logger defines the methods for logger.
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}
