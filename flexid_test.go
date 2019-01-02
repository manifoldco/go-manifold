package manifold

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/manifoldco/go-manifold/idtype"
)

var (
	validID               ID
	validInternalID       *InternalID
	validFlexID           *FlexID
	validFlexIDString     = "web.com" + pathSeperator + "user" + pathSeperator + "abc123"
	validFlexIDJSONString = `"` + strings.Replace(validFlexIDString, "\\", "\\\\", -1) + `"`
	validFlexIDJSONArray  = `["` + strings.Replace(validFlexIDString, "\\", "\", \"", -1) + `"]`
	expectedString        string
	expectedJSONString    string
)

func init() {
	var err error
	validID, err = NewID(idtype.Partner)
	if err != nil {
		panic(err)
	}
	validInternalID = validID.AsComposite()
	validFlexID = validInternalID.AsFlexID().AsFlexID()

	expectedString = fmt.Sprintf("%s%spartner%s%s", ManifoldDomain, pathSeperator,
		pathSeperator, validID)
	expectedJSONString = `"` + strings.Replace(expectedString, "\\", "\\\\", -1) + `"`
}

func TestDomain_Validate(t *testing.T) {
	tests := []struct {
		name    string
		d       Domain
		wantErr bool
	}{
		{
			name:    "Valid Domain",
			d:       "test.com",
			wantErr: false,
		},
		{
			name:    "Invalid Domain",
			d:       "*@#loootbox",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.Validate(nil); (err != nil) != tt.wantErr {
				t.Errorf("Domain.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExternalID_Validate(t *testing.T) {
	tests := []struct {
		name    string
		eid     ExternalID
		wantErr bool
	}{
		{
			name:    "Valid ExternalID",
			eid:     "abc123",
			wantErr: false,
		},
		{
			name:    "Invalid ExternalID",
			eid:     "$$$GETMONEY$$$",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.eid.Validate(nil); (err != nil) != tt.wantErr {
				t.Errorf("ExternalID.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInternalID_Validate(t *testing.T) {
	err := validInternalID.Validate(nil)
	if err != nil {
		t.Errorf("InternalID.Validate() unexpected error: %v", err)
		return
	}
}

func TestFlexID_Validate(t *testing.T) {
	err := validFlexID.Validate(nil)
	if err != nil {
		t.Errorf("FlexID.Validate() unexpected error: %v", err)
		return
	}
}

func TestInternalID_MarshalText(t *testing.T) {
	out, err := validInternalID.MarshalText()
	if err != nil {
		t.Errorf("InternalID.MarshalText() unexpected error: %v", err)
		return
	}
	if string(out) != expectedString {
		t.Errorf("InternalID.MarshalText() expected '%s', got '%s'", expectedString, out)
		return
	}
	var mid *InternalID
	defer func() {
		_ = recover()
	}()
	_, err = mid.MarshalText()
	t.Errorf("InternalID.MarshalText() expected panic when called on nil, got '%s'", err)
}

func TestInternalID_MarshalJSON(t *testing.T) {

	out, err := validInternalID.MarshalJSON()
	if err != nil {
		t.Errorf("InternalID.MarshalJSON() unexpected error: %v", err)
		return
	}
	if string(out) != expectedJSONString {
		t.Errorf("InternalID.MarshalJSON() expected '%s', got '%s'", expectedJSONString, out)
		return
	}

	out, err = json.Marshal(validInternalID)
	if err != nil {
		t.Errorf("json.Marshal(*InternalID) unexpected error: %v", err)
		return
	}
	if string(out) != expectedJSONString {
		t.Errorf("json.Marshal(*InternalID) expected '%s', got '%s'", expectedJSONString, out)
		return
	}

	out, err = json.Marshal(*validInternalID)
	if err != nil {
		t.Errorf("json.Marshal(InternalID) unexpected error: %v", err)
		return
	}
	if string(out) != expectedJSONString {
		t.Errorf("json.Marshal(InternalID) expected '%s', got '%s'", expectedJSONString, out)
		return
	}

	var mid *InternalID
	defer func() {
		_ = recover()
	}()
	_, err = mid.MarshalJSON()
	t.Errorf("InternalID.MarshalJSON() expected panic when called on nil, got '%s'", err)
}

func TestFlexID_MarshalText(t *testing.T) {
	out, err := validFlexID.MarshalText()
	if err != nil {
		t.Errorf("FlexID.MarshalText() unexpected error: %v", err)
		return
	}
	if string(out) != expectedString {
		t.Errorf("FlexID.MarshalText() expected '%s', got '%s'", expectedString, out)
		return
	}
	var fid *FlexID
	defer func() {
		_ = recover()
	}()
	_, err = fid.MarshalText()
	t.Errorf("FlexID.MarshalText() expected panic when called on nil, got '%s'", err)
}

func TestFlexID_MarshalJSON(t *testing.T) {

	out, err := validFlexID.MarshalJSON()
	if err != nil {
		t.Errorf("FlexID.MarshalJSON() unexpected error: %v", err)
		return
	}
	if string(out) != expectedJSONString {
		t.Errorf("FlexID.MarshalJSON() expected '%s', got '%s'", expectedJSONString, out)
		return
	}

	out, err = json.Marshal(validFlexID)
	if err != nil {
		t.Errorf("json.Marshal(*FlexID) unexpected error: %v", err)
		return
	}
	if string(out) != expectedJSONString {
		t.Errorf("json.Marshal(*FlexID) expected '%s', got '%s'", expectedJSONString, out)
		return
	}

	out, err = json.Marshal(*validFlexID)
	if err != nil {
		t.Errorf("json.Marshal(FlexID) unexpected error: %v", err)
		return
	}
	if string(out) != expectedJSONString {
		t.Errorf("json.Marshal(FlexID) expected '%s', got '%s'", expectedJSONString, out)
		return
	}

	var fid *FlexID
	defer func() {
		_ = recover()
	}()
	_, err = fid.MarshalJSON()
	t.Errorf("FlexID.MarshalJSON() expected panic when called on nil, got '%s'", err)
}

func TestInternalID_UnmarshalText(t *testing.T) {
	tests := []struct {
		name     string
		m        *InternalID
		arg      string
		err      error
		expected *InternalID
	}{
		{
			name: "Fails to unmarshal from text when nil",
			arg:  expectedString,
			err:  errNilValue,
		},
		{
			name:     "Unmarshals from text to expected ID",
			m:        &InternalID{},
			arg:      expectedString,
			expected: validInternalID,
		},
		{
			name: "Errors with invalid FlexID",
			m:    &InternalID{},
			arg:  "THIS_IS_TOTALLY_INVALID",
			err:  errInvalidParts,
		},
		{
			name: "Errors with valid FlexID that is not a manifoldID",
			m:    &InternalID{},
			arg:  validFlexIDString,
			err:  ErrNotAInternalID,
		},
		{
			name: "Errors with valid FlexID that is not a manifoldID because of ID",
			m:    &InternalID{},
			arg:  string(ManifoldDomain) + `\user\abc123`,
			err:  ErrNotAInternalID,
		},
		{
			name: "Errors with valid FlexID that is not a manifoldID because of Type mismatch",
			m:    &InternalID{},
			arg:  string(ManifoldDomain) + `\user\` + ID(*validInternalID).String(),
			err:  ErrInternalIDTypeMismatch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.UnmarshalText([]byte(tt.arg)); err != tt.err {
				t.Errorf("InternalID.UnmarshalText() error: %v, expected: %v", err, tt.err)
			} else if tt.expected != nil && *tt.m != *tt.expected {
				t.Errorf("InternalID.UnmarshalText() expected: %v, to equal: %v", tt.m, tt.expected)
			}
		})
	}
}

func TestInternalID_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		m        *InternalID
		arg      string
		err      error
		expected *InternalID
	}{
		{
			name: "Fails to unmarshal from JSON when nil",
			arg:  expectedJSONString,
			err:  errNilValue,
		},
		{
			name:     "Unmarshals from JSON to expected ID",
			m:        &InternalID{},
			arg:      expectedJSONString,
			expected: validInternalID,
		},
		{
			name: "Errors with valid FlexID JSON string that is not a manifoldID",
			m:    &InternalID{},
			arg:  validFlexIDJSONString,
			err:  ErrNotAInternalID,
		},
		{
			name: "Errors with valid FlexID JSON array that is not a manifoldID",
			m:    &InternalID{},
			arg:  validFlexIDJSONArray,
			err:  ErrNotAInternalID,
		},
		{
			name: "Errors with invalid InternalID JSON",
			m:    &InternalID{},
			arg:  "THIS_IS_TOTALLY_INVALID",
			err:  errInvalidParts,
		},
		{
			name: "Errors with valid FlexID JSON string that is not a manifoldID because of ID",
			m:    &InternalID{},
			arg:  `"` + string(ManifoldDomain) + `\\user\\abc123"`,
			err:  ErrNotAInternalID,
		},
		{
			name: "Errors with valid FlexID JSON string that is not a manifoldID because of Type mismatch",
			m:    &InternalID{},
			arg:  `"` + string(ManifoldDomain) + `\\user\\` + ID(*validInternalID).String() + `"`,
			err:  ErrInternalIDTypeMismatch,
		},
		{
			name: "Errors with valid FlexID JSON array that is not a manifoldID because of ID",
			m:    &InternalID{},
			arg:  `["` + string(ManifoldDomain) + `", "user", "abc123"]`,
			err:  ErrNotAInternalID,
		},
		{
			name: "Errors with valid FlexID JSON array that is not a manifoldID because of Type mismatch",
			m:    &InternalID{},
			arg:  `["` + string(ManifoldDomain) + `", "user", "` + ID(*validInternalID).String() + `"]`,
			err:  ErrInternalIDTypeMismatch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.UnmarshalJSON([]byte(tt.arg)); err != tt.err {
				t.Errorf("InternalID.UnmarshalJSON() error: %v, expected: %v", err, tt.err)
			} else if tt.expected != nil && *tt.m != *tt.expected {
				t.Errorf("InternalID.UnmarshalJSON() expected: %v, to equal: %v", tt.m, tt.expected)
			}
		})
	}
}

func TestFlexID_UnmarshalText(t *testing.T) {
	tests := []struct {
		name     string
		id       *FlexID
		arg      string
		err      error
		expected *FlexID
	}{
		{
			name: "Fails to unmarshal from text when nil",
			arg:  expectedString,
			err:  errNilValue,
		},
		{
			name:     "Unmarshals from text to expected ID",
			id:       &FlexID{},
			arg:      expectedString,
			expected: validFlexID,
		},
		{
			name: "Passes with valid FlexID that is not a manifoldID",
			id:   &FlexID{},
			arg:  validFlexIDString,
		},
		{
			name: "Errors with invalid FlexID",
			id:   &FlexID{},
			arg:  "THIS_IS_TOTALLY_INVALID",
			err:  errInvalidParts,
		},
		{
			name: "Errors with invalid FlexID because of Domain",
			id:   &FlexID{},
			arg:  `nope\user\abc123`,
			err:  errInvalidDomain,
		},
		{
			name: "Errors with invalid FlexID because of Type",
			id:   &FlexID{},
			arg:  `nope.com\$$$\abc123`,
			err:  errInvalidType,
		},
		{
			name: "Errors with invalid FlexID because of ID",
			id:   &FlexID{},
			arg:  `nope.com\user\`,
			err:  errInvalidID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.id.UnmarshalText([]byte(tt.arg)); err != tt.err {
				t.Errorf("FlexID.UnmarshalText() error: %v, expected: %v", err, tt.err)
			} else if tt.expected != nil && *tt.id != *tt.expected {
				t.Errorf("FlexID.UnmarshalText() expected: %v, to equal: %v", tt.id, tt.expected)
			}
		})
	}
}

func TestFlexID_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		id       *FlexID
		arg      string
		err      error
		expected *FlexID
	}{
		{
			name: "Fails to unmarshal from JSON when nil",
			arg:  expectedJSONString,
			err:  errNilValue,
		},
		{
			name:     "Unmarshals from JSON to expected ID",
			id:       &FlexID{},
			arg:      expectedJSONString,
			expected: validFlexID,
		},
		{
			name: "Passes with valid FlexID JSON string that is not a manifoldID",
			id:   &FlexID{},
			arg:  validFlexIDJSONString,
		},
		{
			name: "Passes with valid FlexID JSON array that is not a manifoldID",
			id:   &FlexID{},
			arg:  validFlexIDJSONArray,
		},
		{
			name: "Errors with invalid FlexID JSON",
			id:   &FlexID{},
			arg:  "THIS_IS_TOTALLY_INVALID",
			err:  errInvalidParts,
		},
		{
			name: "Errors with invalid FlexID JSON string because of Domain",
			id:   &FlexID{},
			arg:  `"nope\\user\\abc123"`,
			err:  errInvalidDomain,
		},
		{
			name: "Errors with invalid FlexID JSON string because of Type",
			id:   &FlexID{},
			arg:  `"nope.com\\$$$\\abc123"`,
			err:  errInvalidType,
		},
		{
			name: "Errors with invalid FlexID JSON string because of ID",
			id:   &FlexID{},
			arg:  `"nope.com\\user\\"`,
			err:  errInvalidID,
		},
		{
			name: "Errors with invalid FlexID JSON array because of Domain",
			id:   &FlexID{},
			arg:  `["nope", "user", "abc123"]`,
			err:  errInvalidDomain,
		},
		{
			name: "Errors with invalid FlexID JSON array because of Type",
			id:   &FlexID{},
			arg:  `["nope.com", "$$$", "abc123"]`,
			err:  errInvalidType,
		},
		{
			name: "Errors with invalid FlexID JSON array because of ID",
			id:   &FlexID{},
			arg:  `["nope.com", "user", ""]`,
			err:  errInvalidID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.id.UnmarshalJSON([]byte(tt.arg)); err != tt.err {
				t.Errorf("FlexID.UnmarshalJSON() error: %v, expected: %v", err, tt.err)
			} else if tt.expected != nil && *tt.id != *tt.expected {
				t.Errorf("FlexID.UnmarshalJSON() expected: %v, to equal: %v", tt.id, tt.expected)
			}
		})
	}
}

func TestInternalID_AsID(t *testing.T) {
	converted := validInternalID.AsID()
	if converted == nil {
		t.Errorf("InternalID.AsID() got: nil, expected: %v", validID)
	}
	var mid *InternalID
	converted = mid.AsID()
	if converted != nil {
		t.Errorf("InternalID.AsID() got: %s, expected: nil", converted)
	}
}

func TestFlexID_AsInternalID(t *testing.T) {
	converted, err := validFlexID.AsInternalID()
	if err != nil {
		t.Errorf("FlexID.AsInternalID() unexpected error: %v", err)
		return
	}
	if *converted != *validInternalID {
		t.Errorf("FlexID.AsInternalID() got: %v, expected: %v", err, validInternalID)
	}
}

func TestDomain_SubDomain(t *testing.T) {
	tests := []struct {
		name string
		d    Domain
		want string
	}{
		{
			name: "Subdomain returns as expected",
			d:    "test.maniford.co",
			want: "test",
		},
		{
			name: "Multi-segment subdomain returns as expected",
			d:    "test.tony.maniford.co",
			want: "test.tony",
		},
		{
			name: "No subdomain returns empty string",
			d:    ManifoldDomain,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.SubDomain(); got != tt.want {
				t.Errorf("Domain.SubDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
