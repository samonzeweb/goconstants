package goconstants

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Metadata constains metadata of a typed constant.
// Always restrict visibility of variables of Metadata type as
// they are not imutable.
//
// Strings and JSONStrings allows to separated user facing and encoding strings,
// at least one of them must be set.
type Metadata[T comparable] struct {
	// Name of the constant type (used for error messages)
	Name string
	// Strings allow the mapping between a constant value and its string
	// representation.
	// All valid valid must be present in the map as it's used to check
	// the validity of a constant value.
	Strings map[T]string
	// JSONStrings allow the mapping between a constant value and its json
	// representation.
	// All valid valid must be present in the map.
	// If not set, the content of Strings will be used.
	JSONStrings map[T]string
}

// Errors returned by Validate
var (
	ErrNameMissing        = errors.New("the Name field is blank")
	ErrNoStringsDefined   = errors.New("neither Strings not JSONStrings are defined")
	ErrStringsIncoherence = errors.New("Strings and JSONStrings does not have the same keys")
)

// Validate checks that the Metadata instance is valid.
// Use it preferably in a dedicated test.
func (meta Metadata[T]) Validate() error {
	if meta.Name == "" {
		return ErrNameMissing
	}

	if (meta.Strings == nil || len(meta.Strings) == 0) &&
		(meta.JSONStrings == nil || len(meta.JSONStrings) == 0) {
		return ErrNoStringsDefined
	}

	if meta.Strings != nil && meta.JSONStrings != nil {
		if len(meta.Strings) != len(meta.JSONStrings) {
			return ErrStringsIncoherence
		}

		for k := range meta.Strings {
			if _, ok := meta.JSONStrings[k]; !ok {
				return ErrStringsIncoherence
			}
		}
	}

	return nil
}

// StringHelper returns a string representing the constant value.
// If the value is unknown, the string is a blank one.
func (meta Metadata[T]) StringHelper(v T) string {
	s, err := meta.ToStringHelper(v)
	if err != nil {
		return ""
	}

	return s
}

// ToStringHelper returns a string representing the constant value or an error
// if the value is not known. Use it rather than StringHelper if you need to
// check the validity of the value.
func (meta Metadata[T]) ToStringHelper(v T) (string, error) {
	return meta.toStringHelper(v, meta.getStrings())
}

// FromStringHelper converts a string to its associated constant value, and a
// boolean indicating if the value is valid.
// If the boolean is false, ignore the returned value.
func (meta Metadata[T]) FromStringHelper(representation string) (T, bool) {
	return meta.fromStringHelper(representation, meta.getStrings())
}

// toStringHelper returns a string representing the constant value or an error
// if the value is not known.
func (meta Metadata[T]) toStringHelper(v T, strings map[T]string) (string, error) {
	if s, ok := strings[v]; ok {
		return s, nil
	}

	return "", fmt.Errorf("invalid %s value: %#v", meta.Name, v)
}

// fromStringHelper converts a string to its associated constant value, and a
// boolean indicating if the value is valid.
func (meta Metadata[T]) fromStringHelper(representation string, strings map[T]string) (T, bool) {
	for k, v := range strings {
		if representation == v {
			return k, true
		}
	}

	var zero T
	return zero, false
}

// IsValidHelper checks if a given constant is valid (known).
func (meta Metadata[T]) IsValidHelper(v T) bool {
	_, ok := meta.getStrings()[v]
	return ok
}

//getStrings returns Strings if defined (not nil) or JSONStrings
// as fallback.
func (meta Metadata[T]) getStrings() map[T]string {
	if meta.Strings != nil {
		return meta.Strings
	}

	return meta.JSONStrings
}

//getJSONStrings returns JSONStrings if defined (not nil) or Strings
// as fallback.
func (meta Metadata[T]) getJSONStrings() map[T]string {
	if meta.JSONStrings != nil {
		return meta.JSONStrings
	}

	return meta.Strings
}

// MarshalJSONHelper allows the implementation of MarshalJSON for
// the associated constant type.
func (meta Metadata[T]) MarshalJSONHelper(v T) ([]byte, error) {
	representation, err := meta.toStringHelper(v, meta.getJSONStrings())
	if err != nil {
		return nil, fmt.Errorf("unable to mashal %s type to json: %w", meta.Name, err)
	}

	return json.Marshal(representation)
}

// UnmarshalJSONHelper allows the implementation of UnmarshalJSON for
// the associated constant type.
func (meta Metadata[T]) UnmarshalJSONHelper(b []byte, v *T) error {
	var representation string
	err := json.Unmarshal(b, &representation)
	if err != nil {
		return fmt.Errorf("unable to unmashal %s type from json: %w", meta.Name, err)
	}

	value, ok := meta.fromStringHelper(representation, meta.getJSONStrings())
	if !ok {
		return fmt.Errorf("unable to unmashal json, unknown value: %s", representation)
	}

	*v = value
	return nil
}
