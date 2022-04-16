package goconstants_test

import (
	"encoding/json"
	"fmt"

	"github.com/samonzeweb/goconstants"
)

// GopherState is a constant type
type GopherState int

// All valid values for the constant, always use them to set, compare, ...
// Only use iota if the numeric values are never used outside (DB, API, ...)
// as inserting / reordering constants will change the values !
const (
	GopherAsleep GopherState = 1 + iota
	GopherEating
	GopherCoding
	GopherJokingAboutJS
)

// metaGopherState are metadata of GopherState, its visibility is limited.
var metaGopherState = goconstants.Metadata[GopherState]{
	Name: "GopherState",
	Strings: map[GopherState]string{
		GopherAsleep:        "Zzz",
		GopherEating:        "Miam",
		GopherCoding:        "Coding",
		GopherJokingAboutJS: "Lol",
	},
	// JSONStrings is optionnal and not used in the example code.
	// Here how you can enable different strings for JSON content :
	// JSONStrings: map[GopherState]string{
	// 	GopherAsleep:        "sleeping",
	// 	GopherEating:        "eating",
	// 	GopherCoding:        "coding",
	// 	GopherJokingAboutJS: "joking",
	// },
}

// String returns a string representation of the constant.
// It implements the fmt.Stringer interface.
func (gs GopherState) String() string {
	return metaGopherState.StringHelper(gs)
}

// ToString returns string representation of the constant and
// a, error if the given value is unknown.
func (gs GopherState) ToString() (string, error) {
	return metaGopherState.ToStringHelper(gs)
}

// IsValid checks if the given constant is valid (known).
// Useful for validation.
func (gs GopherState) IsValid() bool {
	return metaGopherState.IsValidHelper(gs)
}

// MarshalJSON implement json.Marshaler
func (gs GopherState) MarshalJSON() ([]byte, error) {
	return metaGopherState.MarshalJSONHelper(gs)
}

// UnmarshalJSON implements json.Unmarshaler
func (gs *GopherState) UnmarshalJSON(b []byte) error {
	return metaGopherState.UnmarshalJSONHelper(b, gs)
}

func Example() {
	// Convert to / from strings
	// Using a valid constant value.
	happyGopher := GopherJokingAboutJS
	fmt.Println(happyGopher)
	fmt.Println(happyGopher.ToString())
	fmt.Println(happyGopher.IsValid())

	// Using an invalid constant value (sadly, you can do this).
	var invalidGopher GopherState = 999
	fmt.Println(invalidGopher)
	fmt.Println(invalidGopher.ToString())
	fmt.Println(invalidGopher.IsValid())

	// Parsing valid and invalid constant representations
	fmt.Println(metaGopherState.FromStringHelper("Coding"))
	fmt.Println(metaGopherState.FromStringHelper("Walking on the moon"))

	// Convert to / from JSON
	type Gopher struct {
		Name  string      `json:"name"`
		State GopherState `json:"state"`
	}

	gopher := Gopher{
		Name:  "Sleepy Gopher",
		State: GopherAsleep,
	}

	rawJson, err := json.Marshal(gopher)
	fmt.Println(string(rawJson))
	fmt.Println(err)

	rawJson = []byte(`
		{
			"name":"Geek Gopher",
			"state":"Coding"
		}
	`)

	err = json.Unmarshal(rawJson, &gopher)
	fmt.Printf("%#v\n", gopher)
	fmt.Println(err)

	// Output:
	// Lol
	// Lol <nil>
	// true
	//
	//  invalid GopherState value: 999
	// false
	// Coding true
	//  false
	// {"name":"Sleepy Gopher","state":"Zzz"}
	// <nil>
	// goconstants_test.Gopher{Name:"Geek Gopher", State:3}
	// <nil>
}
