package manifold

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/manifoldco/go-manifold/idtype"
)

var (
	validID               ID
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
	// Call AsFlexID for coverage on both types :D
	validFlexID = validID.AsFlexID().AsFlexID()

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
	err := validID.Validate(nil)
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
			name: "Passes with valid FlexID that is not a Manifold ID",
			id:   &FlexID{},
			arg:  validID.String(),
		},
		{
			name: "Passes with a raw Manifold ID string",
			id:   &FlexID{},
			arg:  validFlexIDString,
		},
		{
			name: "Errors with invalid FlexID with only one part",
			id:   &FlexID{},
			arg:  "THIS_IS_TOTALLY_INVALID",
			err:  errInvalidParts,
		},
		{
			name: "Errors with invalid FlexID with two parts",
			id:   &FlexID{},
			arg:  `THIS_IS_TOTALLY_INVALID\durr`,
			err:  errInvalidParts,
		},
		{
			name: "Errors with invalid FlexID because of Domain",
			id:   &FlexID{},
			arg:  `nope\user\abc123`,
			err:  errInvalidDomain,
		},
		{
			name: "Errors with invalid FlexID because of Class",
			id:   &FlexID{},
			arg:  `nope.com\$$$\abc123`,
			err:  errInvalidClass,
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
			name: "Passes with valid FlexID JSON string that is not a Manifold ID",
			id:   &FlexID{},
			arg:  validFlexIDJSONString,
		},
		{
			name: "Passes with valid FlexID JSON array that is not a Manifold ID",
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
			name: "Errors with invalid FlexID JSON string because of Class",
			id:   &FlexID{},
			arg:  `"nope.com\\$$$\\abc123"`,
			err:  errInvalidClass,
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
			name: "Errors with invalid FlexID JSON array because of Class",
			id:   &FlexID{},
			arg:  `["nope.com", "$$$", "abc123"]`,
			err:  errInvalidClass,
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

func TestFlexID_AsManifoldID(t *testing.T) {
	tests := []struct {
		name     string
		id       FlexID
		err      error
		expected *ID
	}{
		{
			name:     "Succeeds with valid FlexID that is a Manifold ID",
			id:       *validFlexID,
			expected: &validID,
		},
		{
			name: "Errors with valid FlexID that is not a Manifold ID because of Domain",
			id:   FlexID{"doge.co", "user", "abc123"},
			err:  ErrNotAManifoldID,
		},
		{
			name: "Errors with valid FlexID that is not a Manifold ID because of ID",
			id:   FlexID{ManifoldDomain.String(), "user", "abc123"},
			err:  ErrNotAManifoldID,
		},
		{
			name: "Errors with valid FlexID that is not a Manifold ID because of Type mismatch",
			id:   FlexID{ManifoldDomain.String(), "user", validID.String()},
			err:  ErrManifoldIDTypeMismatch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := tt.id.AsManifoldID()
			if err != tt.err {
				t.Errorf("FlexID.AsManifoldID() error: %v, expected: %v", err, tt.err)
			} else if tt.expected != nil && *id != *tt.expected {
				t.Errorf("FlexID.AsManifoldID() expected: %v, to equal: %v", id, *tt.expected)
			}
		})
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

func TestNewFlexID(t *testing.T) {
	tests := []struct {
		name    string
		d       Domain
		c       Class
		id      ExternalID
		want    *FlexID
		wantErr bool
	}{
		{
			name: "Succeeds with valid ID parts",
			d:    validFlexID.Domain(),
			c:    validFlexID.Class(),
			id:   validFlexID.ID(),
			want: validFlexID,
		},
		{
			name:    "Fails with bad Domain",
			d:       Domain("NOTDOMAIN"),
			c:       validFlexID.Class(),
			id:      validFlexID.ID(),
			wantErr: true,
		},
		{
			name:    "Fails with bad Class",
			d:       validFlexID.Domain(),
			c:       Class("%#@$#$#@)"),
			id:      validFlexID.ID(),
			wantErr: true,
		},
		{
			name:    "Fails with bad ID",
			d:       validFlexID.Domain(),
			c:       validFlexID.Class(),
			id:      ExternalID("!@!#@$#*!GETMONEY(^)%#!$(#!$!)("),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFlexID(tt.d, tt.c, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFlexID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFlexID() = %v, want %v", got, tt.want)
			}
		})
	}
}
