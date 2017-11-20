package primitives

import "log"

// Credential represents the specification that is required to filter out
// specific credentials in the Resource spec.
type Credential struct {
	Key     string `json:"key"`
	Name    string `json:"name,omitempty"`
	Default string `json:"default,omitempty"`
}

// Valid will validate the CredentialSpec.
func (c *Credential) Valid() bool {
	if c.Key == "" {
		log.Printf("Credential: invalid key")
		return false
	}

	return true
}

// CredentialValue is a simple representation of the actual key/value of a
// Credential.
type CredentialValue struct {
	Credential `json:",inline"`
	Value      string `json:"value"`
}
