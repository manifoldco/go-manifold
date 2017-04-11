package errors

import "testing"

func TestErrorTypeForStatusCode(t *testing.T) {
	t.Run("with an unknown status code", func(t *testing.T) {
		if et, ok := TypeForStatusCode(1); ok {
			t.Errorf("Expected status code 1 not to return ErrorType `%s`", et)
		}
	})

	t.Run("with a known status code", func(t *testing.T) {
		for key, value := range statusCodeMap {
			if et, ok := TypeForStatusCode(value); !ok && et != key {
				t.Errorf("Expected status code 500 to return `%s`, got `%s`", key, et)
			}
		}
	})
}
