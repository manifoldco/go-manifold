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

func TestType_Collection(t *testing.T) {
	tcs := []struct {
		scenario string
		typ      func() idtype.Type
		plural   string
	}{
		{
			scenario: "user",
			typ: func() idtype.Type {
				return idtype.User
			},
			plural: "users",
		},
		{
			scenario: "resource",
			typ: func() idtype.Type {
				return idtype.Resource
			},
			plural: "resources",
		},
		{
			scenario: "access suffix",
			typ: func() idtype.Type {
				typ := idtype.Type(idtype.TypeOverflow - 2)
				idtype.Register(typ, false, "things_access")

				return typ
			},
			plural: "things_access",
		},
		{
			scenario: "define plurals",
			typ: func() idtype.Type {
				typ := idtype.Type(idtype.TypeOverflow - 1)
				idtype.Register(typ, false, "category", idtype.WithPlural("categories"))

				return typ
			},
			plural: "categories",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			typ := tc.typ()

			if col := typ.Collection(); col != tc.plural {
				t.Errorf("collection '%s' does not match plural '%s'", col, tc.plural)
			}
		})
	}
}

func Test_Register(t *testing.T) {
	tests := map[string]struct {
		fn        func()
		wantPanic bool
	}{
		"simple registration": {
			fn: func() {
				typ := idtype.Type(idtype.TypeOverflow - 2)
				idtype.Register(typ, false, "things_access")
			},
		},
		"double matching registrations": {
			fn: func() {
				typ := idtype.Type(idtype.TypeOverflow - 2)
				idtype.Register(typ, false, "things_access")
				idtype.Register(typ, false, "things_access")
			},
		},
		"double registrations, different names": {
			fn: func() {
				typ := idtype.Type(idtype.TypeOverflow - 2)
				idtype.Register(typ, false, "things_access")
				idtype.Register(typ, false, "something_else")
			},
			wantPanic: true,
		},

		"double registrations, different options": {
			fn: func() {
				typ := idtype.Type(idtype.TypeOverflow - 2)
				idtype.Register(typ, false, "things_access")
				idtype.Register(typ, false, "things_access", idtype.WithPlural("things_accesses"))
			},
			wantPanic: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.wantPanic && r == nil {
					t.Error("Registration should have panicked")
				}
				if !tc.wantPanic && r != nil {
					t.Errorf("Expected success, got panic %v", r)
				}
			}()

			tc.fn()
		})
	}
}
