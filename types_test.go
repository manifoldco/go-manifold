package manifold

import (
	"encoding/json"
	"testing"
)

func TestFeatureMap(t *testing.T) {
	t.Run("different FeatureMaps are not equal", func(t *testing.T) {
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
			t.Error("A and B are not equal, but are.")
		}
		if b.Equals(c) {
			t.Error("B and C are not equal, but are.")
		}
		if a.Equals(c) {
			t.Error("A and C are not equal, but are.")
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
			t.Error("A and B are equal, but aren't.")
		}
	})
}

func TestMetadata(t *testing.T) {
	t.Run("different Metadata are not equal", func(t *testing.T) {
		a := Metadata{
			"abra":       1,
			"bulbasaur":  "TWO",
			"charmander": false,
		}
		b := Metadata{
			"abra":       2,
			"bulbasaur":  "one",
			"charmander": true,
		}
		c := Metadata{
			"abra":      1,
			"bulbasaur": "TWO",
		}

		if a.Equals(b) {
			t.Error("A and B are not equal, but are.")
		}
		if b.Equals(c) {
			t.Error("B and C are not equal, but are.")
		}
		if a.Equals(c) {
			t.Error("A and C are not equal, but are.")
		}
	})

	t.Run("Metadata with the same values are equal", func(t *testing.T) {
		a := Metadata{
			"abra":       1,
			"bulbasaur":  "TWO",
			"charmander": false,
		}
		b := Metadata{
			"abra":       1,
			"bulbasaur":  "TWO",
			"charmander": false,
		}

		if !a.Equals(b) {
			t.Error("A and B are equal, but aren't.")
		}
	})

	t.Run("Expected Metadata is considered valid", func(t *testing.T) {
		mds := []Metadata{
			{
				"abra":       1,
				"bulbasaur":  "TWO",
				"charmander": false,
			},
			{
				"abra":       1,
				"bulbasaur":  "TWO",
				"charmander": false,
				"subdata": Metadata{
					"abra":       1,
					"bulbasaur":  "TWO",
					"charmander": false,
				},
			},
			{
				"abra":       1,
				"bulbasaur":  "TWO",
				"charmander": false,
				"subdata": map[string]interface{}{
					"abra":       1,
					"bulbasaur":  "TWO",
					"charmander": false,
				},
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
				"$$$BADLABEL$$": 1,
			},
			{
				"bad-value-type": json.Encoder{},
			},
			{
				"nested-value-err": Metadata{
					"$$$BADLABEL$$": 1,
				},
			},
			{
				"nested-value-err": map[string]interface{}{
					"$$$BADLABEL$$": 1,
				},
			},
			{
				"json-too-big": makeBigString(MetadataMaxSize),
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
