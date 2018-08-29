package manifold

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

func dataMapEquals(a map[string]interface{}, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if val, ok := b[k]; !ok || val != v {
			return false
		}
	}

	return true
}

// FeatureMap stores the selected feature values for a Manifold resource
type FeatureMap map[string]interface{}

// Equals checks the equality of another FeatureMap against this one
func (f FeatureMap) Equals(fm FeatureMap) bool {
	return dataMapEquals(f, fm)
}

// Metadata stores Metadata for a Manifold resource
type Metadata map[string]interface{}

// MetadataMaxSize defines the max size of the metadata JSON in bytes
const MetadataMaxSize = 1024

// Equals checks the equality of another Metadata against this one
func (m Metadata) Equals(md Metadata) bool {
	return dataMapEquals(m, md)
}

// Validate validates this Metadata
func (m Metadata) Validate(_ interface{}) error {
	for k, v := range m {
		// Make sure key is a label
		if err := Label(k).Validate(nil); err != nil {
			return errors.Errorf("Key '%s' is not a valid Manifold Label", k)
		}
		// Make sure value is an allowed type
		switch t := v.(type) {
		case string, bool, float32, float64, int, int8, int16, int32, int64, uint,
			uint8, uint16, uint32, uint64, json.Number:
			continue
		case map[string]interface{}:
			if err := Metadata(t).Validate(nil); err != nil {
				return err
			}
		case Metadata:
			if err := t.Validate(nil); err != nil {
				return err
			}
		default:
			typ := reflect.TypeOf(v)
			return errors.Errorf("'%s' is an invalid field value type for Metadata", typ.Name())
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
