package manifold

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/dchest/blake2b"

	"github.com/manifoldco/go-base32"

	"github.com/manifoldco/go-manifold/idtype"
)

const (
	idVersion  = 0x1
	byteLength = 18
)

// Comparator for empty IDs
var emptyID [byteLength]byte

// Identifiable is the interface implemented by objects that can be given
// IDs.
type Identifiable interface {
	GetID() ID
	Version() int
	Type() idtype.Type
}

// Immutable structs are Identifiables that do not change, and should be signed.
type Immutable interface {
	Identifiable
	GetBody() interface{}
	Immutable() // We don't ever need to call this, its just for type checking.
}

// Mutable structs are Identifiables that can be changed.
type Mutable interface {
	Identifiable
	Mutable() // also just for type checking.
}

// ID is an encoded unique identifier for an object.
//
// The first byte holds the schema version of the id itself.
// The second byte holds the type of the object.
// The remaining 16 bytes hold a digest of the contents of the object for
// immutable objects, or a random value for mutable objects.
type ID [byteLength]byte

// NewMutableID returns a new ID for a mutable object.
func NewMutableID(body Mutable) (ID, error) {
	t := body.Type()
	return NewID(t)
}

// NewFakeMutableID returns an ID for a fake mutable object, not relying on
// the Body contents of the supplied mutable to generate the ID
func NewFakeMutableID(body Mutable, source string) (ID, error) {
	h, err := blake2b.New(&blake2b.Config{Size: 16})
	if err != nil {
		return ID{}, err
	}
	h.Write([]byte(source))
	preamble := body.Type()
	id := ID{idVersion<<4 | preamble.Upper(), preamble.Lower()}

	copy(id[2:], h.Sum(nil))

	return id, nil
}

// DeriveMutableID returns a ID for a mutable object based on another ID.
func DeriveMutableID(body Mutable, base ID, derivableType idtype.Type) ID {
	preamble := body.Type()
	id := ID{idVersion<<4 | preamble.Upper(), preamble.Lower(), derivableType.Upper(), derivableType.Lower()}
	copy(id[3:], base[3:17])

	return id
}

// NewID returns a new ID for a Mutable idtype using only the Type
func NewID(t idtype.Type) (ID, error) {
	if !t.Mutable() {
		return ID{}, errors.New("Cannot generate ID for non-mutable type")
	}

	id := ID{idVersion<<4 | t.Upper(), t.Lower()}
	_, err := rand.Read(id[2:])
	if err != nil {
		return ID{}, err
	}

	return id, nil
}

// NewImmutableID returns a new signed ID for an immutable object.
//
// sig should be a registry.Signature type
func NewImmutableID(obj Immutable, sig interface{}) (ID, error) {
	h, err := blake2b.New(&blake2b.Config{Size: 16})
	if err != nil {
		return ID{}, err
	}

	h.Write([]byte(strconv.Itoa(obj.Version())))

	b, err := json.Marshal(obj.GetBody())
	if err != nil {
		return ID{}, err
	}
	h.Write(b)

	b, err = json.Marshal(sig)
	if err != nil {
		return ID{}, err
	}
	h.Write(b)

	preamble := obj.Type()
	id := ID{idVersion<<4 | preamble.Upper(), preamble.Lower()}

	copy(id[2:], h.Sum(nil))

	return id, nil
}

// DecodeIDFromString returns an ID that is stored in the given string.
func DecodeIDFromString(value string) (ID, error) {
	buf, err := decodeFromByte([]byte(value))
	if err != nil {
		return ID{}, err
	}

	id := ID{}
	copy(id[:], buf)
	return id, nil
}

// Type returns the idtype.Type encoded object type represented by this ID.
func (id ID) Type() idtype.Type {
	return idtype.Decode(id[0]&0x0F, id[1])
}

func (id ID) String() string {
	return base32.EncodeToString(id[:])
}

// MarshalText implements the encoding.TextMarshaler interface for IDs.
//
// IDs are encoded in unpadded base32.
func (id ID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for IDs.
func (id *ID) UnmarshalText(b []byte) error {
	return id.fillID(b)
}

func (id *ID) fillID(raw []byte) error {
	out, err := decodeFromByte(raw)
	if err != nil {
		return err
	}

	copy(id[:], out)
	return nil
}

func decodeFromByte(raw []byte) ([]byte, error) {
	out, err := base32.DecodeString(string(raw))
	if err != nil {
		return nil, err
	}
	if len(out) != byteLength {
		return nil, errors.New("Incorrect length for id")
	}

	return out, nil
}

// Validate implements the Validate interface for goswagger.
// We know that if the value has successfully parsed, it is valid, so no action
// is required.
func (id ID) Validate(_ interface{}) error {
	return nil
}

// IsEmpty returns whether or not the ID is empty (all zeros)
func (id ID) IsEmpty() bool {
	return id == emptyID
}

// AsComposite converts the ID to a CompositeID
func (id ID) AsComposite() *ManifoldID {
	mid := ManifoldID(id)
	return &mid
}
