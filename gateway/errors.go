package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// Error represents the error returned by a handler to the requestor
type Error struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Class   string `json:"class"`
	Message string `json:"message"`
}

// Error implements the error interface
func (e Error) Error() string {
	return fmt.Sprintf("%s (%s): %s", e.Code, e.Class, e.Message)
}

// ErrorList contains several errors
type ErrorList []Error

// Error implements the error interface
func (el ErrorList) Error() string {
	var sb strings.Builder
	for i, e := range el {
		sb.WriteString(fmt.Sprintf("%s (%s): %s", e.Code, e.Class, e.Message))
		if i < len(el)-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func bodyToErr(body io.Reader) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return errors.New("could not read error body")
	}

	var eList ErrorList
	if err := json.Unmarshal(b, &eList); err == nil {
		return eList
	}

	var e Error
	if err := json.Unmarshal(b, &e); err == nil {
		return e
	}

	return fmt.Errorf("Not a valid error type: %s", b)
}
