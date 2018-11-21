package manifold

import (
	"encoding/json"
	"strings"

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
	Type  MetadataValueType `json:"type"`
	Value interface{}       `json:"value"`
}

// MetadataValueType defines metadata type identifiers
type MetadataValueType string

const (
	// MetadataValueTypeString identifies the string type
	MetadataValueTypeString MetadataValueType = "string"
	// MetadataValueTypeBool identifies the bool type
	MetadataValueTypeBool MetadataValueType = "bool"
	// MetadataValueTypeInt identifies the int type
	MetadataValueTypeInt MetadataValueType = "int"
	// MetadataValueTypeFloat identifies the float type
	MetadataValueTypeFloat MetadataValueType = "float"
	// MetadataValueTypeObject identifies the object type
	MetadataValueTypeObject MetadataValueType = "object"
)

// Equals checks the equality of another MetadataValue against this one
func (m *MetadataValue) Equals(md MetadataValue) bool {
	return m.Type == md.Type && m.Value == md.Value
}

func (m *MetadataValue) tryCastFields() error {
	switch m.Type {
	case MetadataValueTypeString:
		_, ok := m.Value.(string)
		if !ok {
			return errors.New("Expected value to be a string but it was not")
		}
	case MetadataValueTypeBool:
		_, ok := m.Value.(bool)
		if !ok {
			return errors.New("Expected value to be a boolean but it was not")
		}
	case MetadataValueTypeInt:
		val, err := number.ToInt64(m.Value)
		if err != nil {
			return errors.Errorf(
				"Expected value to be castable to int64 but it was not: %s", err.Error())
		}
		m.Value = val
	case MetadataValueTypeFloat:
		val, err := number.ToFloat64(m.Value)
		if err != nil {
			return errors.Errorf(
				"Expected value to be castable to float64 but it was not: %s", err.Error())
		}
		m.Value = val
	case MetadataValueTypeObject:
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
			"%s is not a valid type, expected 'string', 'int', 'float', 'bool', or 'object'", m.Type)
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
	m.Type = MetadataValueType(typ)
	m.Value = val
	if err := m.tryCastFields(); err != nil {
		return errors.Errorf("Could not make MetadataValue from map: %s", err.Error())
	}
	return nil
}

// UnmarshalJSON controls how a MetadataValue is parsed from JSON
func (m *MetadataValue) UnmarshalJSON(data []byte) error {
	// New replica of struct for unmarshal to avoid infinite unmarshal loop
	mv := &struct {
		Type  string      `json:"type"`
		Value interface{} `json:"value"`
	}{}
	json.Unmarshal(data, &mv)
	m.Type = MetadataValueType(mv.Type)
	m.Value = mv.Value
	return m.tryCastFields()
}

// Validate validates this MetadataValue
func (m *MetadataValue) Validate(_ interface{}) error {
	switch m.Type {
	case MetadataValueTypeString:
		_, ok := m.Value.(string)
		if !ok {
			return errors.New("Expected value to be a string but it was not")
		}
	case MetadataValueTypeInt:
		_, ok := m.Value.(int64)
		if !ok {
			return errors.New("Expected value to be a int64 but it was not")
		}
	case MetadataValueTypeFloat:
		_, ok := m.Value.(float64)
		if !ok {
			return errors.New("Expected value to be a float64 but it was not")
		}
	case MetadataValueTypeBool:
		_, ok := m.Value.(bool)
		if !ok {
			return errors.New("Expected value to be a bool but it was not")
		}
	case MetadataValueTypeObject:
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
			"%s is not a valid type, expected 'string', 'int', 'float', 'bool', or 'object'", m.Type)
	}
	return nil
}

// Metadata stores Metadata for a Manifold resource
type Metadata map[string]MetadataValue

// MetadataMaxSize defines the max size of the metadata JSON in bytes
const MetadataMaxSize = 10 * 1024

// ErrMetadataNonexistantKey describes and error for when the expected key is not present
var ErrMetadataNonexistantKey = errors.New("Key does not exist")

// ErrMetadataUnexpectedValueType describes and error when a metadata type isn't what's expected
var ErrMetadataUnexpectedValueType = errors.New("Found value but it was not the expected type")

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

// GetString returns the value of the specified key as a string, or returns an error
func (m Metadata) GetString(key string) (*string, error) {
	val, ok := m[key]
	if !ok {
		return nil, ErrMetadataNonexistantKey
	}
	if val.Type != MetadataValueTypeString {
		return nil, ErrMetadataUnexpectedValueType
	}
	out, _ := val.Value.(string)
	return &out, nil
}

// GetBool returns the value of the specified key as a bool, or returns an error
func (m Metadata) GetBool(key string) (*bool, error) {
	val, ok := m[key]
	if !ok {
		return nil, ErrMetadataNonexistantKey
	}
	if val.Type != MetadataValueTypeBool {
		return nil, ErrMetadataUnexpectedValueType
	}
	out, _ := val.Value.(bool)
	return &out, nil
}

// GetInt returns the value of the specified key as a int64, or returns an error
func (m Metadata) GetInt(key string) (*int64, error) {
	val, ok := m[key]
	if !ok {
		return nil, ErrMetadataNonexistantKey
	}
	if val.Type != MetadataValueTypeInt {
		return nil, ErrMetadataUnexpectedValueType
	}
	out, _ := val.Value.(int64)
	return &out, nil
}

// GetFloat returns the value of the specified key as a float64, or returns an error
func (m Metadata) GetFloat(key string) (*float64, error) {
	val, ok := m[key]
	if !ok {
		return nil, ErrMetadataNonexistantKey
	}
	if val.Type != MetadataValueTypeFloat {
		return nil, ErrMetadataUnexpectedValueType
	}
	out, _ := val.Value.(float64)
	return &out, nil
}

// GetObject returns the value of the specified key as a Metadata, or returns an error
func (m Metadata) GetObject(key string) (Metadata, error) {
	val, ok := m[key]
	if !ok {
		return nil, ErrMetadataNonexistantKey
	}
	if val.Type != MetadataValueTypeObject {
		return nil, ErrMetadataUnexpectedValueType
	}
	out, _ := val.Value.(Metadata)
	return out, nil
}

// AnnotationsMap defines a map of string arrays that contain the annotations data
type AnnotationsMap map[string][]string

// AnnotationMaxReservedKeys defines the max number of reserved keys (Keys prefixed with manifold.co)
const AnnotationMaxReservedKeys = 20

// AnnotationReservedKeyPrefix is the prefix a key must start with to be considered reserved
const AnnotationReservedKeyPrefix = "manifold.co"

// AnnotationKnownReservedKeys is an array of all the known reserved keys, any other key prefixed with the reserved
// key prefix will cause an error.
var AnnotationKnownReservedKeys = []string{
	"manifold.co/tool",
	"manifold.co/package",
	"manifold.co/environment",
}

// Equals checks the equality of another AnnotationsMap against this one
func (a AnnotationsMap) Equals(fm AnnotationsMap) bool {
	if len(a) != len(fm) {
		return false
	}
	for key, value := range a {
		val, ok := fm[key]
		if !ok || len(value) != len(val) {
			return false
		}
		for subkey, subvalue := range value {
			if subval := fm[key][subkey]; subvalue != subval {
				return false
			}
		}
	}

	return true
}

// Validate validates this AnnotationsMap
func (a AnnotationsMap) Validate(_ interface{}) error {
	countReserved := 0
	for key, value := range a {
		// Make sure the key is a valid key
		if err := AnnotationKey(key).Validate(nil); err != nil {
			return errors.Errorf("Key '%s' is not a valid annotation key", key)
		}
		// Make sure that, if the key is reserved, it is validated
		if strings.HasPrefix(key, AnnotationReservedKeyPrefix) {
			found := false
			for _, reservedKey := range AnnotationKnownReservedKeys {
				if reservedKey == key {
					found = true
				}
			}
			if !found {
				return errors.Errorf("Key '%s' is not an accepted reserved key", key)
			}
			countReserved++
		}

		// Make sure every value is a valid value
		for _, subvalue := range value {
			if err := AnnotationValue(subvalue).Validate(nil); err != nil {
				return errors.Errorf("Value '%s' is not a valid annotation value", subvalue)
			}
		}
	}

	// Finally, make sure we didn't overflow the mx number of reserved keys
	if countReserved > AnnotationMaxReservedKeys {
		return errors.New("Annotation has more than 20 annotation keys")
	}
	return nil
}
