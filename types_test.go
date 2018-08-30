package manifold

import (
	json "encoding/json"
	"testing"
)

func TestFeatureMap(t *testing.T) {
	t.Run("Different FeatureMaps are not equal", func(t *testing.T) {
		a := FeatureMap{
			"a": 1,
			"b": "TWO",
			"c": false,
		}
		b := FeatureMap{
			"a": 2,
			"b": "one",
			"c": true,
		}
		c := FeatureMap{
			"a": 1,
			"b": "TWO",
		}

		if a.Equals(b) {
			t.Error("Expected A to not eqaul B, but it did.")
		}
		if b.Equals(c) {
			t.Error("Expected B to not eqaul C, but it did.")
		}
		if a.Equals(c) {
			t.Error("Expected A to not eqaul C, but it did.")
		}
	})

	t.Run("FeatureMaps with the same values are equal", func(t *testing.T) {
		a := FeatureMap{
			"a": 1,
			"b": "TWO",
			"c": false,
		}
		b := FeatureMap{
			"a": 1,
			"b": "TWO",
			"c": false,
		}

		if !a.Equals(b) {
			t.Error("Expected A to eqaul B, but it didn't.")
		}
	})
}

func TestMetadata(t *testing.T) {
	t.Run("Different Metadata are not equal", func(t *testing.T) {
		a := Metadata{
			"abra":       {Type: "int", Value: 1},
			"bulbasaur":  {Type: "string", Value: "TWO"},
			"charmander": {Type: "bool", Value: false},
		}
		b := Metadata{
			"abra":       {Type: "int", Value: 2},
			"bulbasaur":  {Type: "string", Value: "one"},
			"charmander": {Type: "bool", Value: true},
		}
		c := Metadata{
			"abra":      {Type: "int", Value: 1},
			"bulbasaur": {Type: "string", Value: "TWO"},
		}

		if a.Equals(b) {
			t.Error("Expected A to not eqaul B, but it did.")
		}
		if b.Equals(c) {
			t.Error("Expected B to not eqaul C, but it did.")
		}
		if a.Equals(c) {
			t.Error("Expected A to not eqaul C, but it did.")
		}
	})

	t.Run("Metadata with the same values are equal", func(t *testing.T) {
		a := Metadata{
			"abra":       {Type: "int", Value: 1},
			"bulbasaur":  {Type: "string", Value: "TWO"},
			"charmander": {Type: "bool", Value: false},
		}
		b := Metadata{
			"abra":       {Type: "int", Value: 1},
			"bulbasaur":  {Type: "string", Value: "TWO"},
			"charmander": {Type: "bool", Value: false},
		}

		if !a.Equals(b) {
			t.Error("Expected A to eqaul B, but it didn't.")
		}
	})

	t.Run("Expected Metadata is considered valid", func(t *testing.T) {
		mds := []Metadata{
			{
				"abra":       {Type: "int", Value: int64(1)},
				"bulbasaur":  {Type: "string", Value: "TWO"},
				"charmander": {Type: "bool", Value: false},
			},
			{
				"abra":       {Type: "int", Value: int64(1)},
				"bulbasaur":  {Type: "string", Value: "TWO"},
				"charmander": {Type: "bool", Value: false},
				"subdata": {Type: "object", Value: Metadata{
					"abra":       {Type: "int", Value: int64(1)},
					"bulbasaur":  {Type: "string", Value: "TWO"},
					"charmander": {Type: "bool", Value: false},
				}},
			},
		}

		for _, m := range mds {
			if err := m.Validate(nil); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("Expected Metadata is considered invalid", func(t *testing.T) {
		mds := []Metadata{
			{
				"$$$BADLABEL$$": {Type: "int", Value: int64(1)},
			},
			{
				"bad-value-type": {Type: "int", Value: "NOTANINT"},
			},
			{
				"nested-value-err": {Type: "object", Value: Metadata{
					"$$$BADLABEL$$": {Type: "int", Value: int64(1)},
				}},
			},
			{
				"json-too-big": {Type: "string", Value: makeBigString(MetadataMaxSize)},
			},
		}

		for _, m := range mds {
			if err := m.Validate(nil); err == nil {
				t.Errorf("Expected error for Metadata case: %v", m)
			}
		}
	})
}

func makeBigString(size int) string {
	b := make([]byte, size)
	for i := 0; i < size; i++ {
		b[i] = 'a'
	}
	return string(b)
}

func TestMetadataValue(t *testing.T) {

	t.Run("UnmarshalJSON is successful when expected", func(t *testing.T) {
		JSONs := []string{
			`{ "type": "string", "value": "this is a string" }`,
			`{ "type": "bool", "value": false }`,
			`{ "type": "int", "value": 12 }`,
			`{ "type": "float", "value": 12.2 }`,
			`{ "type": "object", "value": {  
				"sub-object-key": { "type": "int", "value": 7 }
			 } }`,
		}

		for _, j := range JSONs {
			mdv := &MetadataValue{}
			if err := json.Unmarshal([]byte(j), &mdv); err != nil {
				t.Error(err)
			}
			if err := mdv.Validate(nil); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("UnmarshalJSON errors when expected", func(t *testing.T) {
		JSONs := []string{
			`{ "type": "string", "value": false }`,
			`{ "type": "bool", "value": 1 }`,
			`{ "type": "int", "value": false }`,
			`{ "type": "float", "value": false }`,
			`{ "type": "object", "value": false }`,
			`{ "type": "snake", "value": false }`,
			`{ "type": "object", "value": {  
				"sub-object-key": false
			 } }`,
			`{ "type": "object", "value": {  
				"sub-object-key": { "value": false }
			 } }`,
			`{ "type": "object", "value": {  
				"sub-object-key": { "type": false, "value": 7 }
			 } }`,
			`{ "type": "object", "value": {  
				"sub-object-key": { "type": "object" }
			 } }`,
			`{ "type": "object", "value": {  
				"sub-object-key": { "type": "snake", "value": false }
			 } }`,
		}

		for _, j := range JSONs {
			mdv := &MetadataValue{}
			if err := json.Unmarshal([]byte(j), &mdv); err == nil {
				t.Errorf("Expected err for Unmarshal of: %s", j)
			}
			if err := mdv.Validate(nil); err == nil {
				t.Errorf("Expected err for Validate of: %s", j)
			}
		}
	})
}
