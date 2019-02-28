package gateway

import (
	"fmt"
	"strings"
)

// gatewayError represents an Error returned by this Middleware to a
// requestor
type gatewayError struct {
	Type     string   `json:"type"`
	Message string `json:"message"`
}

// Error returns the error message represented by this Error
func (e *gatewayError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, strings.Join(e.Messages, ","))
}
