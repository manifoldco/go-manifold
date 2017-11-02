package names

import (
	"testing"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/idtype"
)

func TestNames(t *testing.T) {
	t.Run("two ids get different values", func(t *testing.T) {
		id1, err := manifold.NewID(idtype.Resource)
		if err != nil {
			t.Fatal("Couldn't generate ID")
		}
		id2, err := manifold.NewID(idtype.Resource)
		if err != nil {
			t.Fatal("Couldn't generate ID")
		}

		n1, l1 := New(id1)
		n2, l2 := New(id2)

		if n1 == n2 || l1 == l2 {
			t.Error("generated names match")
		}
	})

	t.Run("the same id gets the same value", func(t *testing.T) {
		id1, err := manifold.NewID(idtype.Resource)
		if err != nil {
			t.Fatal("Couldn't generate ID")
		}

		n1, l1 := New(id1)
		n2, l2 := New(id1)

		if n1 != n2 || l1 != l2 {
			t.Error("generated names don't match")
		}
	})
}
