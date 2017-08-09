package manifold

import (
	"encoding/json"
	"testing"

	"github.com/manifoldco/go-manifold/idtype"
)

func TestIDMarshalUnmarshal(t *testing.T) {
	id, _ := NewID(idtype.User)

	b, err := json.Marshal(id)
	if err != nil {
		t.Fatal("error marshaling", err)
	}

	var out ID
	err = json.Unmarshal(b, &out)
	if err != nil {
		t.Fatal("error unmarshaling", err)
	}

	if id.Equals(out) {
		t.Fatal("did not marshal/unmarshal correctly")
	}
}

func TestIDOmitEmpty(t *testing.T) {
	expected := `{}`

	var in struct {
		Field ID `json:",omitempty"`
	}

	b, err := json.Marshal(&in)
	if err != nil {
		t.Fatal("error marshaling")
	}

	if string(b) != expected {
		t.Fatal("value wasn't omitted. got", string(b))
	}
}

func TestIDUnmarshalZeroLen(t *testing.T) {
	in := `{"Field": "00000000000000000000000000000"}`

	var out struct {
		Field ID
	}

	err := json.Unmarshal([]byte(in), &out)
	if err != nil {
		t.Fatal("error during unmarshal: ", err.Error())
	}

	if len(out.Field) != 0 {
		t.Fatal("did not unmarshal properly")
	}
}
