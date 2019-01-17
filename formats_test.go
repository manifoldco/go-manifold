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

func TestCode(t *testing.T) {
	t.Run("errors on invalid code", func(t *testing.T) {
		c := Code("test")
		err := c.Validate(nil)

		if err == nil {
			t.Error("Expected an error")
		}

		_, ok := err.(*Error)
		if !ok {
			t.Error("Expected a manifold Error")
		}
	})

	t.Run("does not error on valid code", func(t *testing.T) {
		c := Code("0123456789abcdef")
		err := c.Validate(nil)

		if err != nil {
			t.Errorf("Unexpected Error: %s", err)
		}
	})
}

func TestFeatureValueLabel(t *testing.T) {
	t.Run("errors on invalid feature value label", func(t *testing.T) {
		l := FeatureValueLabel("BluesClues")
		err := l.Validate(nil)
		if err == nil {
			t.Error("Expected an error")
		}
		l = FeatureValueLabel("blues clues")
		err = l.Validate(nil)
		if err == nil {
			t.Error("Expected an error")
		}
		l = FeatureValueLabel("moosejuice!")
		err = l.Validate(nil)
		if err == nil {
			t.Error("Expected an error")
		}

		_, ok := err.(*Error)
		if !ok {
			t.Error("Expected a manifold Error")
		}
	})

	t.Run("does not error on valid feature value label", func(t *testing.T) {
		l := FeatureValueLabel("t2.micro")
		err := l.Validate(nil)
		if err != nil {
			t.Errorf("Unexpected Error: %s", err)
		}
		l = FeatureValueLabel("fat-cats")
		err = l.Validate(nil)
		if err != nil {
			t.Errorf("Unexpected Error: %s", err)
		}
		l = FeatureValueLabel("orange_soda")
		err = l.Validate(nil)
		if err != nil {
			t.Errorf("Unexpected Error: %s", err)
		}
	})
}

func TestAnnotationKey(t *testing.T) {
	t.Run("errors on invalid annotation key", func(t *testing.T) {
		l := AnnotationKey("HASMAJ")
		err := l.Validate(nil)
		if err == nil {
			t.Error("Expected an error")
		}
		l = AnnotationKey("/startwithslash")
		err = l.Validate(nil)
		if err == nil {
			t.Error("Expected an error")
		}
		l = AnnotationKey("endswithdot.")
		err = l.Validate(nil)
		if err == nil {
			t.Error("Expected an error")
		}

		_, ok := err.(*Error)
		if !ok {
			t.Error("Expected a manifold Error")
		}
	})

	t.Run("does not error on valid annotation key", func(t *testing.T) {
		l := AnnotationKey("manfiold.co/test")
		err := l.Validate(nil)
		if err != nil {
			t.Errorf("Unexpected Error: %s", err)
		}
		l = AnnotationKey("iam/a.valid-value")
		err = l.Validate(nil)
		if err != nil {
			t.Errorf("Unexpected Error: %s", err)
		}
	})
}

func TestAnnotationValue(t *testing.T) {
	t.Run("errors on invalid annotation value", func(t *testing.T) {
		l := AnnotationValue("has_invalid_char")
		err := l.Validate(nil)
		if err == nil {
			t.Error("Expected an error")
		}
		l = AnnotationValue("/startwithslash")
		err = l.Validate(nil)
		if err == nil {
			t.Error("Expected an error")
		}
		l = AnnotationValue("endswithdot.")
		err = l.Validate(nil)
		if err == nil {
			t.Error("Expected an error")
		}

		_, ok := err.(*Error)
		if !ok {
			t.Error("Expected a manifold Error")
		}
	})

	t.Run("does not error on valid annotation value", func(t *testing.T) {
		l := AnnotationValue("Iam/A.Valid-Value")
		err := l.Validate(nil)
		if err != nil {
			t.Errorf("Unexpected Error: %s", err)
		}
		l = AnnotationValue("alsovalid")
		err = l.Validate(nil)
		if err != nil {
			t.Errorf("Unexpected Error: %s", err)
		}
	})
}

func TestCredentialKey(t *testing.T) {
	t.Run("errors on invalid credential key", func(t *testing.T) {
		l := CredentialKey("abc")
		err := l.Validate(nil)
		if err == nil {
			t.Error("Expected an error")
		}
	})

	t.Run("does not error on valid credential key", func(t *testing.T) {
		l := CredentialKey("ABC")
		err := l.Validate(nil)
		if err != nil {
			t.Errorf("Unexpected Error: %s", err)
		}
	})
}

func TestCredentialBody(t *testing.T) {
	t.Run("errors on invalid credential body", func(t *testing.T) {
		longString := ""

		for len(longString) <= maxCredentialBodySize {
			longString += "a"
		}
		l := CredentialBody(longString)
		err := l.Validate(nil)
		if err == nil {
			t.Error("Expected an error")
		}
	})

	t.Run("does not error on valid credential value", func(t *testing.T) {
		l := CredentialBody("hey, that's not 1024 bytes long!")
		err := l.Validate(nil)
		if err != nil {
			t.Errorf("Unexpected Error: %s", err)
		}
	})
}
