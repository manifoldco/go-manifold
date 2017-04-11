package errors

import "net/http"

// Error Types mapped to their text presentations
var (
	BadRequestError       Type = "bad_request"
	UnauthorizedError     Type = "unauthorized"
	NotFoundError         Type = "not_found"
	ConflictError         Type = "conflict"
	InternalServerError   Type = "internal"
	NotImplementedError   Type = "not_implemented"
	MethodNotAllowedError Type = "method_not_allowed"
)

// Type represents a type of Error returned from the middleware
type Type string

// String returns the string representation of this error type
func (t Type) String() string {
	return string(t)
}

// Code returns the http status code for this type
func (t Type) Code() int {
	return statusCodeMap[t]
}

// TypeForStatusCode returns the Type that is matched for a given
// status code. If no Type is found, it will return false.
func TypeForStatusCode(code int) (Type, bool) {
	// Go doesn't like returning the optional return value immediately so we
	// need to assign it before returning it.
	et, ok := inverseMap[code]
	return et, ok
}

// statusCodeMap relates an error type to an HTTP status code
var statusCodeMap = map[Type]int{
	BadRequestError:       http.StatusBadRequest,
	UnauthorizedError:     http.StatusUnauthorized,
	NotFoundError:         http.StatusNotFound,
	ConflictError:         http.StatusConflict,
	InternalServerError:   http.StatusInternalServerError,
	NotImplementedError:   http.StatusNotImplemented,
	MethodNotAllowedError: http.StatusMethodNotAllowed,
}

var inverseMap = map[int]Type{}

func init() {
	for k, v := range statusCodeMap {
		inverseMap[v] = k
	}
}
