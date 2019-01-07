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
		"Invalid CompositeID, expected 3 parts, Domain, Class, and ID")
	errInvalidDomain = NewError(errors.BadRequestError,
		"Invalid CompositeID, expected a valid Domain in the first segment")
	errInvalidClass = NewError(errors.BadRequestError,
		"Invalid CompositeID, expected a valid Class in the last segment")
	errInvalidID = NewError(errors.BadRequestError,
		"Invalid CompositeID, expected a valid ID in the last segment")

	// ErrNotAManifoldID is an error returned when a CompositeID is expected to
	//  be a Manifold ID, but is not.
	ErrNotAManifoldID = NewError(errors.BadRequestError,
		"Malformed Manifold ID, expected form `manifold.co\\CLASS\\MANIFOLDID`")
	// ErrManifoldIDTypeMismatch is an error returned when a CompositeID is expected to
	//  be a Manifold ID, but is not because the type does not match.
	ErrManifoldIDTypeMismatch = NewError(errors.BadRequestError,
		"Invalid Manifold ID, expected CLASS from `manifold.co\\CLASS\\ID` to match ID Type")
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

// String implements the Stringer interface to easily convert a Domain to a String
func (d Domain) String() string {
	return string(d)
}

// Class is a Manifold Label that is used to represent the class of an ID
type Class Label

// Validate implements the runtime Validatable interface
func (c Class) Validate(r strfmt.Registry) error {
	return c.Label().Validate(r)
}

// String implements the Stringer interface to easily convert a Class to a String
func (c Class) String() string {
	return string(c)
}

// Label implements the Stringer interface to easily convert a Class to a Manifold Label
func (c Class) Label() Label {
	return Label(c)
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

// String implements the Stringer interface to easily convert a ExternalID to a String
func (id ExternalID) String() string {
	return string(id)
}

// CompositeID is an ID that also includes the domain, and type of the identifier.
//  Composed as: DOMAIN / CLASS / ID
//  Example: manifold.co/user/2003btphq7z6dzvjut370jkvkdgcp
//  Has `manifold.co` as the domain, a type of `user`, followed by the Manifold ID.
type CompositeID interface {
	fmt.Stringer
	runtime.Validatable
	encoding.TextMarshaler
	encoding.TextUnmarshaler

	// Domain returns the Domain ( first ) portion of the CompositeID
	Domain() Domain
	// Class returns the Class ( second ) portion of the CompositeID
	Class() Class
	// ID returns the ID ( third ) portion of the CompositeID
	ID() ExternalID
	// AsFlexID allows for easy conversion of all CompositeIDs to the most forgiving struct
	AsFlexID() *FlexID
}

// NewFlexID constructs a FlexID from the provided Domain, Class, and ID parts
func NewFlexID(d Domain, c Class, id ExternalID) (*FlexID, error) {
	out := FlexID{d.String(), c.String(), id.String()}
	if err := out.Validate(nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// FlexIDFromID takes a Manifold ID and converts it to a FlexID for storage
func FlexIDFromID(id ID) *FlexID {
	return &FlexID{id.Domain().String(), id.Class().String(), id.ID().String()}
}

// FlexID is an implementation of CompositeID that is designed to store internal
//  and external IDs it could still store InternalIDs but the InternalID type is
//  preferred as it is directly translatable to a `ID`
type FlexID [3]string

// Domain returns the domain portion as a string
func (id FlexID) Domain() Domain {
	return Domain(id[0])
}

// Class returns the type portion as string
func (id FlexID) Class() Class {
	return Class(id[1])
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
	return fmt.Sprintf("%s%s%s%s%s", id.Domain(), pathSeperator, id.Class(),
		pathSeperator, id.ID())
}

// Validate implements the Validate interface for goswagger
//  which always succeeds because the ID is already parsed
func (id FlexID) Validate(v strfmt.Registry) error {
	if err := id.Domain().Validate(v); err != nil {
		return err
	}
	if id.Class().Validate(v) != nil {
		return errInvalidClass
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
	switch len(parts) {
	case 1:
		// Try to parse as a Manifold ID
		mid := &ID{}
		if err := mid.UnmarshalText(b); err != nil {
			return errInvalidParts
		}
		copy(id[:], (*mid.AsFlexID())[:])
	case 3:
		copy(id[:], parts)
	default:
		return errInvalidParts
	}

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
	// First to attempt to unmarshal as array, though this is undesired as a storage
	//  format it's easy to support for translation
	var parts [3]string
	if err := json.Unmarshal(b, &parts); err != nil {
		// Attempt to unmarshal as string
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return errInvalidParts
		}
		// Leverage the text unmarshalling now we have the string
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

// AsManifoldID validates that the FlexID adheres with the requirements of a ManifoldID
//  and attempts to cast it to one
func (id FlexID) AsManifoldID() (*ID, error) {
	if id.Domain() != ManifoldDomain {
		return nil, ErrNotAManifoldID
	}
	mid, err := DecodeIDFromString(id.ID().String())
	if err != nil {
		return nil, ErrNotAManifoldID
	}
	if mid.Type().Name() != id.Class().String() {
		return nil, ErrManifoldIDTypeMismatch
	}
	return &mid, nil
}

// Ensure interface adherence
var (
	_ runtime.Validatable      = Domain("")
	_ fmt.Stringer             = Domain("")
	_ runtime.Validatable      = Class("")
	_ fmt.Stringer             = Class("")
	_ runtime.Validatable      = ExternalID("")
	_ fmt.Stringer             = ExternalID("")
	_ CompositeID              = &FlexID{}
	_ encoding.TextMarshaler   = &FlexID{}
	_ encoding.TextUnmarshaler = &FlexID{}
	_ json.Marshaler           = &FlexID{}
	_ json.Unmarshaler         = &FlexID{}
	_ CompositeID              = &ID{}
)
