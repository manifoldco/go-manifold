package manifold

import (
	"testing"
)

func TestEmail(t *testing.T) {
	t.Run("errors on invalid email", func(t *testing.T) {
		e := Email("test")
		err := e.Validate(nil)

		if err == nil {
			t.Error("Expected an error")
		}

		_, ok := err.(*Error)
		if !ok {
			t.Error("Expected a manifold Error")
		}
	})

	t.Run("does not error on valid email", func(t *testing.T) {
		e := Email("test@test.com")
		err := e.Validate(nil)

		if err != nil {
			t.Errorf("Unexpected Error: %s", err)
		}
	})
}
