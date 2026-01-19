package id

import (
	"database/sql/driver"

	"github.com/gofrs/uuid/v5"
)

var UUID7Nil = UUID7(uuid.Nil)

// UUID7SentinelMax is a UUID greater than any valid UUID7.
// UUID7 structure:
// | 48 bits timestamp | 4 bits version | 12 bits subsec |
// | 2–3 bits variant  | 62 bits random |
// the version is always 7 (0111) and the variant bits are either 10 or 11,
// so the maximum possible UUID7 is ffffffff-ffff-7fff-bfff-ffffffffffff
var UUID7SentinelMax = UUID7(uuid.Must(uuid.FromString("ffffffff-ffff-7fff-bfff-ffffffffffff")))

// UUID7 represents a UUID version 7.
type UUID7 uuid.UUID

func (u UUID7) String() string {
	return uuid.UUID(u).String()
}

func (u UUID7) Bytes() []byte {
	b := uuid.UUID(u)
	return b[:]
}

// MarshalJSON implements the json.Marshaler interface.
func (u UUID7) MarshalJSON() ([]byte, error) {
	// Encode as a JSON string
	return []byte(`"` + u.String() + `"`), nil
}

// JSONUnmarshaler implements the json.Unmarshaler interface.
func (u *UUID7) UnmarshalJSON(data []byte) error {
	// Expect a JSON string
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return uuid.ErrInvalidFormat
	}

	s := string(data[1 : len(data)-1])

	parsed, err := uuid.FromString(s)
	if err != nil {
		return err
	}

	return assignValidated(u, parsed)
}

// Scan implements the sql.Scanner interface.
func (u *UUID7) Scan(src any) error {
	if src == nil {
		*u = UUID7Nil
		return nil
	}

	switch v := src.(type) {
	case []byte:
		parsed, err := uuid.FromBytes(v)
		if err != nil {
			return err
		}
		return assignValidated(u, parsed)

	case string:
		parsed, err := uuid.FromString(v)
		if err != nil {
			return err
		}
		return assignValidated(u, parsed)

	default:
		return uuid.ErrInvalidFormat
	}
}

// Value implements the driver.Valuer interface.
func (u UUID7) Value() (driver.Value, error) {
	return uuid.UUID(u), nil
}

type UUID7GeneratorParser interface {
	UUID7Generator
	UUID7Parser
}

type UUID7Generator interface {
	Generate() UUID7
}

type UUID7Parser interface {
	ParseString(idStr string) (UUID7, error)
	ParseBytes(idBytes []byte) (UUID7, error)
}

type uuid7Tool struct{}

func NewUUIDV7Tool() *uuid7Tool {
	return &uuid7Tool{}
}

func (g *uuid7Tool) Generate() UUID7 {
	return UUID7(uuid.Must(uuid.NewV7()))
}

// ParseString parses a UUIDv7 from its string representation.
func (g *uuid7Tool) ParseString(idStr string) (UUID7, error) {
	id, err := uuid.FromString(idStr)
	if err != nil {
		return UUID7Nil, err
	}

	var uid UUID7
	uid, err = validate(id)
	if err != nil {
		return UUID7Nil, err
	}

	return uid, nil
}

// ParseBytes parses a UUIDv7 from its byte slice representation.
func (g *uuid7Tool) ParseBytes(idBytes []byte) (UUID7, error) {
	id, err := uuid.FromBytes(idBytes)
	if err != nil {
		return UUID7Nil, err
	}

	var uid UUID7
	uid, err = validate(id)
	if err != nil {
		return UUID7Nil, err
	}

	return uid, nil
}

func validate(u uuid.UUID) (UUID7, error) {
	if u.Version() != uuid.V7 {
		return UUID7Nil, ErrInvalidVersion7
	}
	if u.Variant() != uuid.VariantRFC4122 {
		return UUID7Nil, ErrInvalidVariant
	}
	return UUID7(u), nil
}

func assignValidated(dst *UUID7, u uuid.UUID) error {
	uid, err := validate(u)
	if err != nil {
		return err
	}
	*dst = uid
	return nil
}
