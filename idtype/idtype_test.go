package idtype_test

import (
	"testing"

	"github.com/manifoldco/go-manifold/idtype"
)

func TestUpperAndLower(t *testing.T) {

	if idtype.Product.Upper() != 0x00 {
		t.Error("Decoded upper byte was incorrect")
	}

	if idtype.Product.Lower() != 0xC9 {
		t.Error("Decoded lower byte was incorrect")
	}
}

func TestDecode(t *testing.T) {
	o := idtype.Decode(0x00, 0xC9)

	if o != idtype.Product {
		t.Error("Decoded type did not match Product")
	}
}

func TestTypeFromString(t *testing.T) {
	t.Run("with an existing type", func(t *testing.T) {
		tp := idtype.TypeFromString("user")

		if tp != idtype.User {
			t.Errorf("Expected type to be 'User', got '%s'", tp)
		}
	})

	t.Run("with a non-existing type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected a panic when asking for a non-existing Type")
			}
		}()

		idtype.TypeFromString("non-existing")
	})
}
