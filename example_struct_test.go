package goconstants_test

import (
	"encoding/json"
	"fmt"

	"github.com/samonzeweb/goconstants"
)

// LanguageEnum is an enumeration of programming languages.
// Using a struct with a non exported int lead to a safer constant type
// as it's not possible to affect a value of another type though implicit
// conversion.
// It will be less convenient when used with SQL, and will requires the
// implementation of sql.Scanner and driver.Valuer interfaces.
type LanguageEnum struct {
	value int
}

// All valid values for the LanguageEnum constant.
// The counterpart of using a struct is that we can't use const, and the values
// could be altered.
// If a default value is required, create a constant using the default value
// in the inner value (here 0).
var (
	Go     = LanguageEnum{1}
	Rust   = LanguageEnum{2}
	Python = LanguageEnum{3}
)

// metaGopherState are metadata of GopherState, its visibility is limited.
var metaLanguageEnum = goconstants.Metadata[LanguageEnum]{
	Name: "GopherState",
	Strings: map[LanguageEnum]string{
		Go:     "Go",
		Rust:   "Rust",
		Python: "Python",
	},
}

// String returns a string representation of the constant.
// It implements the fmt.Stringer interface.
func (gs LanguageEnum) String() string {
	return metaLanguageEnum.StringHelper(gs)
}

// ToString returns string representation of the constant and
// a, error if the given value is unknown.
func (gs LanguageEnum) ToString() (string, error) {
	return metaLanguageEnum.ToStringHelper(gs)
}

// IsValid checks if the given constant is valid (known).
// Useful for validation.
func (gs LanguageEnum) IsValid() bool {
	return metaLanguageEnum.IsValidHelper(gs)
}

// MarshalJSON implement json.Marshaler
func (gs LanguageEnum) MarshalJSON() ([]byte, error) {
	return metaLanguageEnum.MarshalJSONHelper(gs)
}

// UnmarshalJSON implements json.Unmarshaler
func (gs *LanguageEnum) UnmarshalJSON(b []byte) error {
	return metaLanguageEnum.UnmarshalJSONHelper(b, gs)
}

func Example_struct() {
	// Convert to / from strings
	// Using a valid constant value.
	language := Rust
	fmt.Println(language)
	fmt.Println(language.ToString())
	fmt.Println(language.IsValid())

	// Using an invalid constant.
	// language = 999 // will not compile
	invalidLanguage := LanguageEnum{} // will compile !
	fmt.Println(invalidLanguage)
	fmt.Println(invalidLanguage.ToString())
	fmt.Println(invalidLanguage.IsValid())

	// Parsing valid and invalid constant representations
	fmt.Println(metaLanguageEnum.FromStringHelper("Rust"))
	fmt.Println(metaLanguageEnum.FromStringHelper("Basic"))

	// Convert to / from JSON
	type SoftwareProduct struct {
		Name     string       `json:"name"`
		Language LanguageEnum `json:"language"`
	}

	software := SoftwareProduct{
		Name:     "Django",
		Language: Python,
	}

	rawJson, err := json.Marshal(software)
	fmt.Println(string(rawJson))
	fmt.Println(err)

	rawJson = []byte(`
		{
			"name":"Hugo",
			"language":"Go"
		}
	`)

	err = json.Unmarshal(rawJson, &software)
	fmt.Printf("%#v\n", software)
	fmt.Println(err)

	// Output:
	// Rust
	// Rust <nil>
	// true
	//
	//  invalid GopherState value: goconstants_test.LanguageEnum{value:0}
	// false
	// Rust true
	//  false
	// {"name":"Django","language":"Python"}
	// <nil>
	// goconstants_test.SoftwareProduct{Name:"Hugo", Language:goconstants_test.LanguageEnum{value:1}}
	// <nil>
}
