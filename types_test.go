package manifold

import (
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
