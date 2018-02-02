package manifold

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
