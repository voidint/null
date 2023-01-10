package null

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Uint8 is an nullable uint8.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Uint8 struct {
	sql.NullByte
}

// NewUint8 creates a new Uint8
func NewUint8(i uint8, valid bool) Uint8 {
	return Uint8{
		NullByte: sql.NullByte{
			Byte:  i,
			Valid: valid,
		},
	}
}

// Uint8From creates a new Uint8 that will always be valid.
func Uint8From(i uint8) Uint8 {
	return NewUint8(i, true)
}

// Uint8FromPtr creates a new Uint8 that be null if i is nil.
func Uint8FromPtr(i *uint8) Uint8 {
	if i == nil {
		return NewUint8(0, false)
	}
	return NewUint8(*i, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Uint8) ValueOrZero() uint8 {
	if !i.Valid {
		return 0
	}
	return i.Byte
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number, string, and null input.
// 0 will not be considered a null Uint8.
func (i *Uint8) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		i.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &i.Byte); err != nil {
		var typeError *json.UnmarshalTypeError
		if errors.As(err, &typeError) {
			// special case: accept string input
			if typeError.Value != "string" {
				return fmt.Errorf("null: JSON input is invalid type (need int or string): %w", err)
			}
			var str string
			if err := json.Unmarshal(data, &str); err != nil {
				return fmt.Errorf("null: couldn't unmarshal number string: %w", err)
			}
			n, err := strconv.ParseUint(str, 10, 8)
			if err != nil {
				return fmt.Errorf("null: couldn't convert string to int: %w", err)
			}
			i.Byte = uint8(n)
			i.Valid = true
			return nil
		}
		return fmt.Errorf("null: couldn't unmarshal JSON: %w", err)
	}

	i.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint8 if the input is blank.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint8) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	n, err := strconv.ParseUint(string(text), 10, 8)
	if err != nil {
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}
	i.Byte = uint8(n)
	i.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Uint8 is null.
func (i Uint8) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(int64(i.Byte), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Uint8 is null.
func (i Uint8) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(i.Byte), 10)), nil
}

// SetValid changes this Uint8's value and also sets it to be non-null.
func (i *Uint8) SetValid(n uint8) {
	i.Byte = n
	i.Valid = true
}

// Ptr returns a pointer to this Uint8's value, or a nil pointer if this Int is null.
func (i Uint8) Ptr() *uint8 {
	if !i.Valid {
		return nil
	}
	return &i.Byte
}

// IsZero returns true for invalid Ints, for future omitempty support (Go 1.4?)
// A non-null Int with a 0 value will not be considered zero.
func (i Uint8) IsZero() bool {
	return !i.Valid
}

// Equal returns true if both ints have the same value or are both null.
func (i Uint8) Equal(other Uint8) bool {
	return i.Valid == other.Valid && (!i.Valid || i.Byte == other.Byte)
}
