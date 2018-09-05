package number

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestToInt64(t *testing.T) {

	tcs := []struct {
		scenario string
		in       interface{}
		out      int64
		err      string
	}{
		{scenario: "int8", in: int8(123), out: 123},
		{scenario: "int16", in: int16(12345), out: 12345},
		{scenario: "int32", in: int32(12345), out: 12345},
		{scenario: "int64", in: int64(12345), out: 12345},
		{scenario: "int", in: 12345, out: 12345},
		{
			scenario: "float32 doesn't lose precision",
			in:       float32(12345),
			out:      12345,
		},
		{
			scenario: "float32 loses precision",
			in:       float32(123.45),
			err:      "invalid casting float32 (123.449997) to int64, precision loss detected",
		},
		{
			scenario: "float64 doesn't lose precision",
			in:       float64(12345),
			out:      12345,
		},
		{
			scenario: "float64 loses precision",
			in:       float64(123.45),
			err:      "invalid casting float64 (123.450000) to int64, precision loss detected",
		},
		{
			scenario: "string",
			in:       "12345",
			out:      12345,
		},
		{
			scenario: "json.Number",
			in:       json.Number("12345"),
			out:      12345,
		},
		{
			scenario: "uint64",
			in:       uint64(12345),
			err:      fmt.Sprintf("Unhandled type: %t", interface{}(uint64(12345))),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			out, err := ToInt64(tc.in)

			if tc.err == "" {
				if out != tc.out {
					t.Errorf("Output %d did not match expected %d", out, tc.out)
				}
			} else {
				if err.Error() != tc.err {
					t.Errorf("Error '%s' did not match expected '%s'", err.Error(), tc.err)
				}
			}
		})
	}

}

func TestToFloat64(t *testing.T) {

	tcs := []struct {
		scenario string
		in       interface{}
		out      float64
		err      string
	}{
		{scenario: "int8", in: int8(123), out: float64(123)},
		{scenario: "int16", in: int16(12345), out: float64(12345)},
		{scenario: "int32", in: int32(12345), out: float64(12345)},
		{scenario: "int64", in: int64(12345), out: float64(12345)},
		{scenario: "int", in: 12345, out: float64(12345)},
		{scenario: "float32", in: float32(12345), out: float64(12345)},
		{scenario: "float64", in: float64(12345), out: float64(12345)},
		{scenario: "string", in: "12345", out: float64(12345)},
		{scenario: "json.Number", in: json.Number("12345.5"), out: float64(12345.5)},
		{
			scenario: "uint64",
			in:       uint64(12345),
			err:      fmt.Sprintf("Unhandled type: %t", interface{}(uint64(12345))),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			out, err := ToFloat64(tc.in)

			if tc.err == "" {
				if out != tc.out {
					t.Errorf("Output %f did not match expected %f", out, tc.out)
				}
			} else {
				if err.Error() != tc.err {
					t.Errorf("Error '%s' did not match expected '%s'", err.Error(), tc.err)
				}
			}
		})
	}

}
