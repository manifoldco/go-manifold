package primitives

import (
	"fmt"
)

// Resource is the specification that is required to build a valid Resource
// manifest.
type Resource struct {
	Name        string        `json:"resource,name"`
	Team        string        `json:"team,omitempty"`
	Credentials []*Credential `json:"credentials,omitempty"`
}

// Valid will validate the Resource.
func (r *Resource) Valid() bool {
	if r.Name == "" {
		fmt.Println("no resource spec label")
		return false
	}

	for _, c := range r.Credentials {
		if !c.Valid() {
			return false
		}
	}

	return true
}
