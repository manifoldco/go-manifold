package manifold

import (
	"encoding/json"

	"github.com/manifoldco/go-manifold/number"
	"github.com/pkg/errors"
)

// FeatureMap stores the selected feature values for a Manifold resource
type FeatureMap map[string]interface{}

// Equals checks the equality of another FeatureMap against this one
func (f FeatureMap) Equals(fm FeatureMap) bool {
	if len(f) != len(fm) {
		return false
	}
	for k, v := range f {
		if val, ok := fm[k]; !ok || val != v {
			return false
		}
	}

	return true
}

// MetadataValue stores MetadataValue for a Manifold resource
type MetadataValue struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// Equals checks the equality of another MetadataValue against this one
func (m *MetadataValue) Equals(md MetadataValue) bool {
	return m.Type == md.Type && m.Value == md.Value
}

func (m *MetadataValue) tryCastFields() error {
	switch m.Type {
	case "string":
		_, ok := m.Value.(string)
		if !ok {
			return errors.New("Expected value to be a string but it was not")
		}
	case "bool":
		_, ok := m.Value.(bool)
		if !ok {
			return errors.New("Expected value to be a boolean but it was not")
		}
	case "int":
		val, err := number.ToInt64(m.Value)
		if err != nil {
			return errors.Errorf(
				"Expected value to be castable to int64 but it was not: %s", err.Error())
		}
		m.Value = val
	case "float":
		val, err := number.ToFloat64(m.Value)
		if err != nil {
			return errors.Errorf(
				"Expected value to be castable to float64 but it was not: %s", err.Error())
		}
		m.Value = val
	case "object":
		val, ok := m.Value.(Metadata)
		if !ok {
			valMap, ok := m.Value.(map[string]interface{})
			if !ok {
				return errors.Errorf("Expected value to be a valid metadata object but it was not")
			}
			val = Metadata{}
			for k, v := range valMap {
				vMap, ok := v.(map[string]interface{})
				if !ok {
					return errors.Errorf(
						"Expected value to be a valid metadata object but it's values aren't compatible")
				}
				kv := MetadataValue{}
				if err := kv.FromMap(vMap); err != nil {
					return errors.Errorf(
						"Expected value to be a valid metadata object but it was not: %s", err.Error())
				}
				val[k] = kv
			}
		}
		m.Value = val
	default:
		return errors.Errorf(
			"%s is not a valid type, expected 'string', 'int', 'float', 'bool', or 'object", m.Type)
	}
	return nil
}

// FromMap tries to make a MetadataValue from a supplied map
func (m *MetadataValue) FromMap(md map[string]interface{}) error {
	typI, ok := md["type"]
	if !ok {
		return errors.New("Could not make MetadataValue from map, it did not contain a 'type' key")
	}
	typ, ok := typI.(string)
	if !ok {
		return errors.New("Could not make MetadataValue from map, the type key was not a string as expected")
	}
	val, ok := md["value"]
	if !ok {
		return errors.New("Could not make MetadataValue from map, it did not contain a 'value' key")
	}
	m.Type = typ
	m.Value = val
	if err := m.tryCastFields(); err != nil {
		return errors.Errorf("Could not make MetadataValue from map: %s", err.Error())
	}
	return nil
}

// UnmarshalJSON controls how a MetadataValue is parsed from JSON
func (m *MetadataValue) UnmarshalJSON(data []byte) error {
	mv := &struct {
		Type  string      `json:"type"`
		Value interface{} `json:"value"`
	}{}
	json.Unmarshal(data, &mv)
	m.Type = mv.Type
	m.Value = mv.Value
	return m.tryCastFields()
}

// Validate validates this MetadataValue
func (m *MetadataValue) Validate(_ interface{}) error {
	switch m.Type {
	case "string":
		_, ok := m.Value.(string)
		if !ok {
			return errors.New("Expected value to be a string but it was not")
		}
	case "int":
		_, ok := m.Value.(int64)
		if !ok {
			return errors.New("Expected value to be a int64 but it was not")
		}
	case "float":
		_, ok := m.Value.(float64)
		if !ok {
			return errors.New("Expected value to be a float64 but it was not")
		}
	case "bool":
		_, ok := m.Value.(bool)
		if !ok {
			return errors.New("Expected value to be a bool but it was not")
		}
	case "object":
		val, ok := m.Value.(Metadata)
		if !ok {
			return errors.New("Expected value to be a Metadata but it was not")
		}
		err := val.Validate(nil)
		if err != nil {
			return errors.Errorf("Metadata value was not valid: %s", err.Error())
		}
	default:
		return errors.Errorf(
			"%s is not a valid type, expected 'string', 'int', 'float', 'bool', or 'object", m.Type)
	}
	return nil
}

// Metadata stores Metadata for a Manifold resource
type Metadata map[string]MetadataValue

// MetadataMaxSize defines the max size of the metadata JSON in bytes
const MetadataMaxSize = 10 * 1024

// Equals checks the equality of another Metadata against this one
func (m Metadata) Equals(md Metadata) bool {
	if len(m) != len(md) {
		return false
	}
	for k, v := range m {
		if val, ok := md[k]; !ok || !val.Equals(v) {
			return false
		}
	}

	return true
}

// Validate validates this Metadata
func (m Metadata) Validate(_ interface{}) error {
	for k, v := range m {
		// Make sure key is a label
		if err := Label(k).Validate(nil); err != nil {
			return errors.Errorf("Key '%s' is not a valid Manifold Label", k)
		}
		// Make sure value is valid
		if err := v.Validate(nil); err != nil {
			return errors.Errorf("Value of key '%s' is not valid: %s", k, err.Error())
		}
	}
	// Make sure total length isn't too long
	b, _ := json.Marshal(m)
	bLen := len(b)
	if bLen > MetadataMaxSize {
		return errors.Errorf("%d is %d bytes larger than the Metadata size limit of %d",
			bLen, MetadataMaxSize-bLen, MetadataMaxSize)
	}
	return nil
}
