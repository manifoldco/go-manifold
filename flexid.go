package manifold

import (
	"encoding"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/manifoldco/go-manifold/errors"
)

const (
	pathSeperator string = `\`

	// ManifoldDomain is the domain name used to identify Manifold IDs
	ManifoldDomain Domain = "manifold.co"
)

var (
	// domainRegex expects that the string is a valid and easy to understand hostname
	domainRegex = regexp.
			MustCompile(`^((?:[a-zA-Z0-9-_]+\.)*)[a-zA-Z0-9][a-zA-Z0-9-_]+\.[a-zA-Z]{2,11}?$`)
	// idRegex expects that an ID at least has a length of one, an only includes
	//  characters expected in Base64 encoded values, GUIDs and UUIDs
	idRegex = regexp.MustCompile(`^\{?[a-zA-Z0-9+/-_]{1,256}={0,2}\}?$`)

	errNilValue = NewError(errors.InternalServerError,
		"Invalid CompositeID, cannot unmarshal to nil ID")
	errInvalidParts = NewError(errors.BadRequestError,
		"Invalid CompositeID, expected 3 parts, Domain, Type, and ID")
	errInvalidDomain = NewError(errors.BadRequestError,
		"Invalid CompositeID, expected a valid Domain in the first segment")
	errInvalidType = NewError(errors.BadRequestError,
		"Invalid CompositeID, expected a valid Type in the last segment")
	errInvalidID = NewError(errors.BadRequestError,
		"Invalid CompositeID, expected a valid ID in the last segment")

	// ErrNotAInternalID is an error returned when a CompositeID is expected to
	//  be a InternalID, but is not.
	ErrNotAInternalID = NewError(errors.BadRequestError,
		"Malformed InternalID, expected form `manifold.co\\TYPE\\MANIFOLDID`")
	// ErrInternalIDTypeMismatch is an error returned when a CompositeID is expected to
	//  be a InternalID, but is not because the type does not match.
	ErrInternalIDTypeMismatch = NewError(errors.BadRequestError,
		"Invalid InternalID, expected TYPE from `manifold.co\\TYPE\\ID` to match ID Type")
)

// Domain is a string that can be Validated based on Regex to expect a string
//  that represents a Domain
type Domain string

// Validate ensures the name value is valid
func (d Domain) Validate(_ strfmt.Registry) error {
	if domainRegex.Match([]byte(d)) {
		return nil
	}

	return errInvalidDomain
}

// SubDomain returns the subdomain portion of the domain
func (d Domain) SubDomain() string {
	parts := domainRegex.FindSubmatch([]byte(d))

	if len(parts[1]) > 0 {
		subdomain := string(parts[1])
		return subdomain[:len(subdomain)-1]
	}

	return ""
}

// ExternalID is a string that can be Validated based on Regex to expect a string
//  representative of an ExternalID
type ExternalID string

// Validate ensures the name value is valid
func (eid ExternalID) Validate(_ strfmt.Registry) error {
	if idRegex.Match([]byte(eid)) {
		return nil
	}

	return errInvalidID
}

// CompositeID is an ID that also includes the domain, and type of the identifier.
//  Composed as: DOMAIN / TYPE / ID
//  Example: manifold.co/user/2003btphq7z6dzvjut370jkvkdgcp
//  Has `manifold.co` as the domain, a type of `user`, followed by the Manifold ID.
type CompositeID interface {
	// Domain returns the Domain ( first ) portion of the CompositeID
	Domain() Domain
	// Type returns the Type ( second ) portion of the CompositeID
	Type() Label
	// ID returns the ID ( third ) portion of the CompositeID
	ID() ExternalID
	// AsFlexID allows for easy conversion of all CompositeIDs to the most forgiving struct
	AsFlexID() *FlexID
	// Stringer interface for easy translation to string
	String() string
	// Validate allows for OpenAPI validation of our structs so we can use them in
	//  OpenAPI schemas
	Validate(strfmt.Registry) error
	// MarshalText allows CompositeIDs to be easily converted to text
	MarshalText() ([]byte, error)
	// UnmarshalText allows CompositeIDs to be easily parsed from text
	UnmarshalText(b []byte) error
	// MarshalText allows CompositeIDs to be easily converted to text
	MarshalJSON() ([]byte, error)
	// UnmarshalText allows CompositeIDs to be easily parsed from text
	UnmarshalJSON(b []byte) error
}

// InternalID is an implementation of CompositeID that wraps the existing Manifold ID type.
//  This allow us to quickly convert existing IDs to the CompositeID format
type InternalID ID

// Domain returns the domain portion as a string
func (m InternalID) Domain() Domain {
	return ManifoldDomain
}

// Type returns the type portion as string
func (m InternalID) Type() Label {
	return Label(ID(m).Type().Name())
}

// ID returns the ID portion as a string
func (m InternalID) ID() ExternalID {
	return ExternalID(ID(m).String())
}

// AsFlexID returns the ID as the FlexID type as required by the CompositeID interface
func (m InternalID) AsFlexID() *FlexID {
	return &FlexID{string(m.Domain()), string(m.Type()), string(m.ID())}
}

// String implements the Stringer interface for go
func (m InternalID) String() string {
	return fmt.Sprintf("%s%s%s%s%s", m.Domain(), pathSeperator, m.Type(),
		pathSeperator, m.ID())
}

// Validate implements the Validate interface for goswagger
func (m InternalID) Validate(v strfmt.Registry) error {
	return ID(m).Validate(v)
}

// MarshalText implements the encoding.TextMarshaler interface
func (m InternalID) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (m *InternalID) UnmarshalText(b []byte) error {
	if m == nil {
		return errNilValue
	}
	id := FlexID{}
	if err := id.UnmarshalText(b); err != nil {
		return err
	}
	mid, err := id.AsInternalID()
	if err != nil {
		return err
	}
	copy(m[:], mid[:])
	return err
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface
func (m *InternalID) UnmarshalJSON(b []byte) error {
	if m == nil {
		return errNilValue
	}
	id := &FlexID{}
	if err := id.UnmarshalJSON(b); err != nil {
		return err
	}
	mid, err := id.AsInternalID()
	if err != nil {
		return err
	}
	copy(m[:], mid[:])
	return err
}

// MarshalJSON implements the encoding/json.Marshaler interface
func (m InternalID) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

// AsID casts the InternalID pointer to a ID pointer for convenience
func (m *InternalID) AsID() *ID {
	if m == nil {
		return nil
	}
	id := ID(*m)
	return &id
}

// FlexID is an implementation of CompositeID that is designed to store internal
//  and external IDs it could still store InternalIDs but the InternalID type is
//  preferred as it is directly translatable to a `ID`
type FlexID [3]string

// Domain returns the domain portion as a string
func (id FlexID) Domain() Domain {
	return Domain(id[0])
}

// Type returns the type portion as string
func (id FlexID) Type() Label {
	return Label(id[1])
}

// ID returns the ID portion as a string
func (id FlexID) ID() ExternalID {
	return ExternalID(id[2])
}

// AsFlexID returns the ID as the FlexID type as required by the CompositeID interface
func (id FlexID) AsFlexID() *FlexID {
	return &id
}

// String implements the Stringer interface for go
func (id FlexID) String() string {
	return fmt.Sprintf("%s%s%s%s%s", id.Domain(), pathSeperator, id.Type(),
		pathSeperator, id.ID())
}

// Validate implements the Validate interface for goswagger
//  which always succeeds because the ID is already parsed
func (id FlexID) Validate(v strfmt.Registry) error {
	if err := id.Domain().Validate(v); err != nil {
		return err
	}
	if id.Type().Validate(v) != nil {
		return errInvalidType
	}
	if err := id.ID().Validate(v); err != nil {
		return err
	}

	return nil
}

// MarshalText implements the encoding.TextMarshaler interface
func (id FlexID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (id *FlexID) UnmarshalText(b []byte) error {
	if id == nil {
		return errNilValue
	}
	parts := strings.Split(string(b), pathSeperator)
	if len(parts) != 3 {
		return errInvalidParts
	}

	copy(id[:], parts)

	if err := id.Validate(nil); err != nil {
		return err
	}

	return nil
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface
func (id *FlexID) UnmarshalJSON(b []byte) error {
	if id == nil {
		return errNilValue
	}
	var parts [3]string
	if err := json.Unmarshal(b, &parts); err != nil {
		// Attempt to unmarshal as string
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return errInvalidParts
		}
		return id.UnmarshalText([]byte(s))
	}

	copy(id[:], parts[:])

	if err := id.Validate(nil); err != nil {
		return err
	}

	return nil
}

// MarshalJSON implements the encoding/json.Marshaler interface
func (id FlexID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

// AsInternalID validates that the FlexID adheres with the requirements of a InternalID
//  and attempts to cast it to one
func (id FlexID) AsInternalID() (*InternalID, error) {
	if id.Domain() != ManifoldDomain {
		return nil, ErrNotAInternalID
	}
	mid, err := DecodeIDFromString(string(id.ID()))
	if err != nil {
		return nil, ErrNotAInternalID
	}
	if mid.Type().Name() != string(id.Type()) {
		return nil, ErrInternalIDTypeMismatch
	}
	out := InternalID(mid)
	return &out, nil
}

// Ensure interface adherence
var (
	_ runtime.Validatable      = Domain("")
	_ runtime.Validatable      = ExternalID("")
	_ CompositeID              = &InternalID{}
	_ fmt.Stringer             = &InternalID{}
	_ runtime.Validatable      = &InternalID{}
	_ encoding.TextMarshaler   = &InternalID{}
	_ encoding.TextUnmarshaler = &InternalID{}
	_ json.Marshaler           = &InternalID{}
	_ json.Unmarshaler         = &InternalID{}
	_ CompositeID              = &FlexID{}
	_ fmt.Stringer             = &FlexID{}
	_ runtime.Validatable      = &FlexID{}
	_ encoding.TextMarshaler   = &FlexID{}
	_ encoding.TextUnmarshaler = &FlexID{}
	_ json.Marshaler           = &FlexID{}
	_ json.Unmarshaler         = &FlexID{}
)