package number

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// ToInt64 takes an interface and converts many things to int64 if able, returns error if it fails
func ToInt64(val interface{}) (int64, error) {
	switch v := val.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float32:
		f := val.(float32)
		i := int64(f)

		if f != float32(i) {
			return 0, errors.Errorf("invalid casting float32 (%f) to int64, precision loss detected", f)
		}
		return i, nil
	case float64:
		f := val.(float64)
		i := int64(f)

		if f != float64(i) {
			return 0, errors.Errorf("invalid casting float64 (%f) to int64, precision loss detected", f)
		}
		return i, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case json.Number:
		return v.Int64()
	default:
		return 0, fmt.Errorf("Unhandled type: %t", v)
	}
}

// ToFloat64 takes an interface and converts many things to float64 if able, returns error if it fails
func ToFloat64(val interface{}) (float64, error) {
	switch v := val.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	case json.Number:
		return v.Float64()
	default:
		return 0, fmt.Errorf("Unhandled type: %t", v)
	}
}
