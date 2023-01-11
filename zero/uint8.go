package zero

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Int is a nullable uint8.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
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

// Uint8From creates a new Uint8 that will be null if zero.
func Uint8From(i uint8) Uint8 {
	return NewUint8(i, i != 0)
}

// Uint8FromPtr creates a new Uint8 that be null if i is nil.
func Uint8FromPtr(i *uint8) Uint8 {
	if i == nil {
		return NewUint8(0, false)
	}
	n := NewUint8(*i, true)
	return n
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Uint8) ValueOrZero() uint8 {
	if !i.Valid {
		return 0
	}
	return i.Byte
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Int.
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
				return fmt.Errorf("zero: JSON input is invalid type (need int or string): %w", err)
			}
			var str string
			if err := json.Unmarshal(data, &str); err != nil {
				return fmt.Errorf("zero: couldn't unmarshal number string: %w", err)
			}
			n, err := strconv.ParseUint(str, 10, 8)
			if err != nil {
				return fmt.Errorf("zero: couldn't convert string to int: %w", err)
			}
			i.Byte = uint8(n)
			i.Valid = n != 0
			return nil
		}
		return fmt.Errorf("zero: couldn't unmarshal JSON: %w", err)
	}

	i.Valid = i.Byte != 0
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint8 if the input is a blank, or zero.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint8) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	n, err := strconv.ParseUint(string(text), 10, 8)
	if err != nil {
		return fmt.Errorf("zero: couldn't unmarshal text: %w", err)
	}
	i.Byte = uint8(n)
	i.Valid = i.Byte != 0
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode 0 if this Int is null.
func (i Uint8) MarshalJSON() ([]byte, error) {
	n := i.Byte
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(int64(n), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Int is null.
func (i Uint8) MarshalText() ([]byte, error) {
	n := i.Byte
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(int64(n), 10)), nil
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

// IsZero returns true for null or zero Ints, for future omitempty support (Go 1.4?)
func (i Uint8) IsZero() bool {
	return !i.Valid || i.Byte == 0
}

// Equal returns true if both ints have the same value or are both either null or zero.
func (i Uint8) Equal(other Uint8) bool {
	return i.ValueOrZero() == other.ValueOrZero()
}
