package connector

import (
	"fmt"
)

// Error is a dappctrl server error.
type Error struct {
	Status  int    `json:"-"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf(
		"server responed with error: %s (%d)", e.Message, e.Code)
}
