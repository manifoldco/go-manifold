package primitives

import (
	"fmt"
)

// Project is the specification that is required to build a valid Project
// manifest.
type Project struct {
	Name      string      `json:"project,name"`
	Team      string      `json:"team,omitempty"`
	Resources []*Resource `json:"resources,omitempty"`
}

// Valid will validate the Project.
func (p *Project) Valid() bool {
	if p.Name == "" {
		fmt.Println("no project spec name")
		return false
	}

	for _, r := range p.Resources {
		if !r.Valid() {
			return false
		}
	}

	return true
}
