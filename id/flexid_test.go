package id

import (
	"fmt"
	"testing"

	"github.com/manifoldco/go-manifold"
	"github.com/manifoldco/go-manifold/idtype"
)

var (
	validManifoldID   *ManifoldID
	validFlexID       *FlexID
	validFlexIDString string = "web.com" + pathSeperator + "user" + pathSeperator + "abc123"
	expectedString    string
)

func init() {
	validID, err := manifold.NewID(idtype.Partner)
	if err != nil {
		panic(err)
	}
	validManifoldID = FromID(validID)
	validFlexID = validManifoldID.AsFlexID().AsFlexID()

	expectedString = fmt.Sprintf("%s%spartner%s%s", ManifoldDomain, pathSeperator,
		pathSeperator, validID)
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

func TestManifoldID_Validate(t *testing.T) {
	err := validManifoldID.Validate(nil)
	if err != nil {
		t.Errorf("ManifoldID.Validate() unexpected error: %v", err)
		return
	}
}

func TestFlexID_Validate(t *testing.T) {
	err := validFlexID.Validate(nil)
	if err != nil {
		t.Errorf("ManifoldID.Validate() unexpected error: %v", err)
		return
	}
}

func TestManifoldID_MarshalText(t *testing.T) {
	out, err := validManifoldID.MarshalText()
	if err != nil {
		t.Errorf("ManifoldID.MarshalText() unexpected error: %v", err)
		return
	}
	if string(out) != expectedString {
		t.Errorf("ManifoldID.MarshalText() expected '%s', got '%s'", expectedString, out)
		return
	}
}

func TestFlexID_MarshalText(t *testing.T) {
	out, err := validFlexID.MarshalText()
	if err != nil {
		t.Errorf("ManifoldID.MarshalText() unexpected error: %v", err)
		return
	}
	if string(out) != expectedString {
		t.Errorf("ManifoldID.MarshalText() expected '%s', got '%s'", expectedString, out)
		return
	}
}

func TestManifoldID_UnmarshalText(t *testing.T) {
	tests := []struct {
		name string
		m    *ManifoldID
		arg  string
		err  error
	}{
		{
			name: "Unmarshals from text when nil",
			arg:  expectedString,
		},
		{
			name: "Unmarshals from text when not nil",
			m:    &ManifoldID{},
			arg:  expectedString,
		},
		{
			name: "Errors with invalid FlexID",
			arg:  "THIS_IS_TOTALLY_INVALID",
			err:  errInvalidParts,
		},
		{
			name: "Errors with valid FlexID that is not a manifoldID",
			arg:  validFlexIDString,
			err:  ErrNotAManifoldID,
		},
		{
			name: "Errors with valid FlexID that is not a manifoldID because of ID",
			arg:  string(ManifoldDomain) + `\user\abc123`,
			err:  ErrNotAManifoldID,
		},
		{
			name: "Errors with valid FlexID that is not a manifoldID because of Type mismatch",
			arg:  string(ManifoldDomain) + `\user\` + manifold.ID(*validManifoldID).String(),
			err:  ErrManifoldIDTypeMismatch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.UnmarshalText([]byte(tt.arg)); err != tt.err {
				t.Errorf("ManifoldID.UnmarshalText() error: %v, expected: %v", err, tt.err)
			}
		})
	}
}

func TestFlexID_UnmarshalText(t *testing.T) {
	tests := []struct {
		name string
		id   *FlexID
		arg  string
		err  error
	}{
		{
			name: "Unmarshals from text when nil",
			arg:  expectedString,
		},
		{
			name: "Unmarshals from text when not nil",
			id:   &FlexID{},
			arg:  expectedString,
		},
		{
			name: "Passes with valid FlexID that is not a manifoldID",
			arg:  validFlexIDString,
		},
		{
			name: "Errors with invalid FlexID",
			arg:  "THIS_IS_TOTALLY_INVALID",
			err:  errInvalidParts,
		},
		{
			name: "Errors with invalid FlexID because of Domain",
			arg:  `nope\user\abc123`,
			err:  errInvalidDomain,
		},
		{
			name: "Errors with invalid FlexID because of Type",
			arg:  `nope.com\$$$\abc123`,
			err:  errInvalidType,
		},
		{
			name: "Errors with invalid FlexID because of ID",
			arg:  `nope.com\user\`,
			err:  errInvalidID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.id.UnmarshalText([]byte(tt.arg)); err != tt.err {
				t.Errorf("FlexID.UnmarshalText() error: %v, expected: %v", err, tt.err)
			}
		})
	}
}
