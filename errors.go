package manifold

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/runtime"

	"github.com/manifoldco/go-manifold/errors"
)

// HTTPError interface represents an error that is returned to a user as an
// HTTP response
type HTTPError interface {
	error
	StatusCode() int
	WriteResponse(http.ResponseWriter, runtime.Producer)
}

// Error represents an Error returned by this Middleware to a
// requestor
type Error struct {
	Type     errors.Type `json:"type"`
	Messages []string    `json:"message"`
}

// NewError returns an Error containing 1 or more error messages
func NewError(t errors.Type, m ...string) *Error {
	return &Error{
		Type:     t,
		Messages: m,
	}
}

// FromError returns an error of type Error from a struct that represents a
func FromError(err error) *Error {
	apiErr, ok := err.(*Error)
	if ok {
		return apiErr
	}

	return NewError(errors.InternalServerError, "Internal Server Error")
}

// Error returns the error message represented by this Error
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type.String(), strings.Join(e.Messages, ","))
}

// Code returns the HTTP Status Code associated with this Error, completes the
// go-openapi error interface.
func (e *Error) Code() int32 {
	return int32(e.Type.Code())
}

// StatusCode returns the HTTP Status Code associated with this error,
// completes the HTTPError interface.
func (e *Error) StatusCode() int {
	return e.Type.Code()
}

// WriteResponse completes the interface for a middleware.Responder from
// go-openapi/runtime
//
// A panic will occur if the given producer errors.
func (e *Error) WriteResponse(w http.ResponseWriter, pr runtime.Producer) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Type.Code())
	if err := pr.Produce(w, e); err != nil {
		panic(err) // recovery middleware will handle it
	}
}

// ToError receives an error and mutates it into an Error based on the concrete
// type of the Error.
func ToError(err error) HTTPError {
	switch e := err.(type) {
	case *Error:
		return e
	case HTTPError:
		if et, ok := errors.TypeForStatusCode(e.StatusCode()); ok {
			return NewError(et, e.Error())
		}
	}

	return NewError(errors.InternalServerError, "Internal Server Error")
}
