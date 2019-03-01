package gateway

import (
	"fmt"
)

// Error represents the error returned by a handler to the requestor
type Error struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Class   string `json:"class"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e *Error) Error() string {
	return fmt.Sprintf("%s (%s): %s", e.Code, e.Class, e.Message)
}
