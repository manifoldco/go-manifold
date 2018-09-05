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
			"abra":       {Type: MetadataValueTypeInt, Value: 1},
			"bulbasaur":  {Type: MetadataValueTypeString, Value: "TWO"},
			"charmander": {Type: MetadataValueTypeBool, Value: false},
		}
		b := Metadata{
			"abra":       {Type: MetadataValueTypeInt, Value: 2},
			"bulbasaur":  {Type: MetadataValueTypeString, Value: "one"},
			"charmander": {Type: MetadataValueTypeBool, Value: true},
		}
		c := Metadata{
			"abra":      {Type: MetadataValueTypeInt, Value: 1},
			"bulbasaur": {Type: MetadataValueTypeString, Value: "TWO"},
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
			"abra":       {Type: MetadataValueTypeInt, Value: 1},
			"bulbasaur":  {Type: MetadataValueTypeString, Value: "TWO"},
			"charmander": {Type: MetadataValueTypeBool, Value: false},
		}
		b := Metadata{
			"abra":       {Type: MetadataValueTypeInt, Value: 1},
			"bulbasaur":  {Type: MetadataValueTypeString, Value: "TWO"},
			"charmander": {Type: MetadataValueTypeBool, Value: false},
		}

		if !a.Equals(b) {
			t.Error("Expected A to eqaul B, but it didn't.")
		}
	})

	t.Run("Expected Metadata is considered valid", func(t *testing.T) {
		mds := []Metadata{
			{
				"abra":       {Type: MetadataValueTypeInt, Value: int64(1)},
				"bulbasaur":  {Type: MetadataValueTypeString, Value: "TWO"},
				"charmander": {Type: MetadataValueTypeBool, Value: false},
			},
			{
				"abra":       {Type: MetadataValueTypeInt, Value: int64(1)},
				"bulbasaur":  {Type: MetadataValueTypeString, Value: "TWO"},
				"charmander": {Type: MetadataValueTypeBool, Value: false},
				"subdata": {Type: MetadataValueTypeObject, Value: Metadata{
					"abra":       {Type: MetadataValueTypeInt, Value: int64(1)},
					"bulbasaur":  {Type: MetadataValueTypeString, Value: "TWO"},
					"charmander": {Type: MetadataValueTypeBool, Value: false},
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
				"$$$BADLABEL$$": {Type: MetadataValueTypeInt, Value: int64(1)},
			},
			{
				"bad-value-type": {Type: MetadataValueTypeInt, Value: "NOTANINT"},
			},
			{
				"nested-value-err": {Type: MetadataValueTypeObject, Value: Metadata{
					"$$$BADLABEL$$": {Type: MetadataValueTypeInt, Value: int64(1)},
				}},
			},
			{
				"json-too-big": {Type: MetadataValueTypeString, Value: makeBigString(MetadataMaxSize)},
			},
		}

		for _, m := range mds {
			if err := m.Validate(nil); err == nil {
				t.Errorf("Expected error for Metadata case: %v", m)
			}
		}
	})

	t.Run("HasType functions, function as expected", func(t *testing.T) {
		m := Metadata{
			"int":    {Type: MetadataValueTypeInt, Value: int64(1)},
			"float":  {Type: MetadataValueTypeFloat, Value: float64(1.23)},
			"string": {Type: MetadataValueTypeString, Value: "izzastring"},
			"bool":   {Type: MetadataValueTypeBool, Value: false},
			"object": {Type: MetadataValueTypeObject, Value: Metadata{
				"abra":    {Type: MetadataValueTypeInt, Value: int64(1)},
				"kadabra": {Type: MetadataValueTypeInt, Value: int64(2)},
			}},
		}
		expectFailFuncs := map[MetadataValueType]func(string){
			MetadataValueTypeInt: func(k string) {
				iVal, err := m.GetInt(k)
				if iVal != nil || err != ErrMetadataUnexpectedValueType {
					t.Errorf("Expected unexpected value type error for Metadata GetInt, got (%v),(%v)", iVal, err)
				}
			},
			MetadataValueTypeFloat: func(k string) {
				fVal, err := m.GetFloat(k)
				if fVal != nil || err != ErrMetadataUnexpectedValueType {
					t.Errorf("Expected unexpected value type error for Metadata GetFloat, got (%v),(%v)", fVal, err)
				}
			},
			MetadataValueTypeString: func(k string) {
				sVal, err := m.GetString(k)
				if sVal != nil || err != ErrMetadataUnexpectedValueType {
					t.Errorf("Expected unexpected value type error for Metadata GetString, got (%v),(%v)", sVal, err)
				}
			},
			MetadataValueTypeBool: func(k string) {
				bVal, err := m.GetBool(k)
				if bVal != nil || err != ErrMetadataUnexpectedValueType {
					t.Errorf("Expected unexpected value type error for Metadata GetBool, got (%v),(%v)", bVal, err)
				}
			},
			MetadataValueTypeObject: func(k string) {
				oVal, err := m.GetObject(k)
				if oVal != nil || err != ErrMetadataUnexpectedValueType {
					t.Errorf("Expected unexpected value type error for Metadata GetObject, got (%v),(%v)", oVal, err)
				}
			},
		}

		for k, v := range m {
			switch v.Type {
			case MetadataValueTypeInt:
				iVal, err := m.GetInt(k)
				if iVal == nil || err != nil {
					t.Errorf("Expected no error for Metadata GetInt, got (%v),(%v)", iVal, err)
				}
			case MetadataValueTypeFloat:
				fVal, err := m.GetFloat(k)
				if fVal == nil || err != nil {
					t.Errorf("Expected no error for Metadata GetFloat, got (%v),(%v)", fVal, err)
				}
			case MetadataValueTypeString:
				sVal, err := m.GetString(k)
				if sVal == nil || err != nil {
					t.Errorf("Expected no error for Metadata GetString, got (%v),(%v)", sVal, err)
				}
			case MetadataValueTypeBool:
				bVal, err := m.GetBool(k)
				if bVal == nil || err != nil {
					t.Errorf("Expected no error for Metadata GetBool, got (%v),(%v)", bVal, err)
				}
			case MetadataValueTypeObject:
				oVal, err := m.GetObject(k)
				if oVal == nil || err != nil {
					t.Errorf("Expected no error for Metadata GetObject, got (%v),(%v)", oVal, err)
				}
			}

			for typ, failTest := range expectFailFuncs {
				if typ == v.Type {
					// Skip, tested earlier for success
					continue
				}
				failTest(k)
			}
		}

	})

	t.Run("HasType functions, fail with missing keys as expected", func(t *testing.T) {
		m := Metadata{}

		iVal, err := m.GetInt("notakey")
		if iVal != nil || err != ErrMetadataNonexistantKey {
			t.Errorf("Expected nonexistant key error for Metadata GetInt, got (%v),(%v)", iVal, err)
		}
		fVal, err := m.GetFloat("notakey")
		if fVal != nil || err != ErrMetadataNonexistantKey {
			t.Errorf("Expected nonexistant key error for Metadata GetFloat, got (%v),(%v)", fVal, err)
		}
		sVal, err := m.GetString("notakey")
		if sVal != nil || err != ErrMetadataNonexistantKey {
			t.Errorf("Expected nonexistant key error for Metadata GetString, got (%v),(%v)", sVal, err)
		}
		bVal, err := m.GetBool("notakey")
		if bVal != nil || err != ErrMetadataNonexistantKey {
			t.Errorf("Expected nonexistant key error for Metadata GetBool, got (%v),(%v)", bVal, err)
		}
		oVal, err := m.GetObject("notakey")
		if oVal != nil || err != ErrMetadataNonexistantKey {
			t.Errorf("Expected nonexistant key error for Metadata GetObject, got (%v),(%v)", oVal, err)
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
				t.Errorf("Failed to Unmarshal (%s) : %s", j, err.Error())
			}
			if err := mdv.Validate(nil); err != nil {
				t.Errorf("Failed to Validate (%v) : %s", mdv, err.Error())
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
