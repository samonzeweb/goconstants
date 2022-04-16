package goconstants_test

import (
	"reflect"
	"testing"

	"github.com/samonzeweb/goconstants"
)

func TestValidate(t *testing.T) {
	type dummy int
	const dummyValueOne dummy = 1
	const dummyValueTwo dummy = 2
	dummyStrings := map[dummy]string{
		dummyValueOne: "I'm a dummy value",
		dummyValueTwo: "I'm another dummy value",
	}
	incompleteStrings := map[dummy]string{
		dummyValueOne: "I'm a dummy value",
	}
	emptyStrings := make(map[dummy]string)

	noNameCase := goconstants.Metadata[dummy]{Strings: dummyStrings}
	err := noNameCase.Validate()
	if err != goconstants.ErrNameMissing {
		t.Errorf("validate didn't catch missing name")
	}

	noStringsCase := goconstants.Metadata[dummy]{Name: "dummy"}
	err = noStringsCase.Validate()
	if err != goconstants.ErrNoStringsDefined {
		t.Errorf("validate didn't catch missing strings")
	}

	emptyStringsCase := goconstants.Metadata[dummy]{
		Name:        "dummy",
		Strings:     emptyStrings,
		JSONStrings: emptyStrings,
	}
	err = emptyStringsCase.Validate()
	if err != goconstants.ErrNoStringsDefined {
		t.Errorf("validate didn't catch empty strings")
	}

	differentStringsLenthCase := goconstants.Metadata[dummy]{
		Name:        "dummy",
		Strings:     dummyStrings,
		JSONStrings: incompleteStrings,
	}
	err = differentStringsLenthCase.Validate()
	if err != goconstants.ErrStringsIncoherence {
		t.Errorf("validate didn't catch different length strings")
	}

	differentStringsKeysCase := goconstants.Metadata[dummy]{
		Name:    "dummy",
		Strings: dummyStrings,
		JSONStrings: map[dummy]string{
			1000: "oups",
			1001: "another oups",
		},
	}
	err = differentStringsKeysCase.Validate()
	if err != goconstants.ErrStringsIncoherence {
		t.Errorf("validate didn't catch different strings keys")
	}
}

type simpson int

const (
	homer  simpson = 1
	marge          = 2
	bart           = 3
	lisa           = 4
	maggie         = 5
)

var cstMeta = goconstants.Metadata[simpson]{
	Name: "cst",
	Strings: map[simpson]string{
		homer:  "Homer Simpson",
		marge:  "Marge Simpson",
		bart:   "Bart Simpson",
		lisa:   "Lisa Simpson",
		maggie: "Maggie Simpson",
	},
	JSONStrings: map[simpson]string{
		homer:  "homer_simpson",
		marge:  "marge_simpson",
		bart:   "bart_simpson",
		lisa:   "lisa_simpson",
		maggie: "maggie_simpson",
	},
}

func TestStringHelper(t *testing.T) {
	testCases := []struct {
		name     string
		input    simpson
		expected string
	}{
		{
			name:     "valid value homer",
			input:    homer,
			expected: "Homer Simpson"},
		{
			name:     "valid value lisa",
			input:    lisa,
			expected: "Lisa Simpson"},
		{
			name:     "invalid value",
			input:    0,
			expected: ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			representation := cstMeta.StringHelper(testCase.input)

			if representation != testCase.expected {
				t.Errorf("expected %s, got %s", testCase.expected, representation)
			}
		})
	}
}

func TestToStringHelper(t *testing.T) {
	testCases := []struct {
		name                string
		input               simpson
		expectedValue       string
		expectedErrorString string
	}{
		{
			name:  "valid value homer",
			input: homer, expectedValue: "Homer Simpson",
			expectedErrorString: "",
		},
		{
			name:  "valid value lisa",
			input: lisa, expectedValue: "Lisa Simpson",
			expectedErrorString: "",
		},
		{
			name:  "invalid value",
			input: 999, expectedValue: "", expectedErrorString: "invalid cst value: 999",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			representation, err := cstMeta.ToStringHelper(testCase.input)

			if representation != testCase.expectedValue {
				t.Errorf("expected value %s, got %s", testCase.expectedValue, representation)
			}

			if err != nil {
				if testCase.expectedErrorString == "" {
					t.Errorf("unexpected error %#v", err)
				} else {
					errorString := err.Error()
					if errorString != testCase.expectedErrorString {
						t.Errorf("expected error %s, got %s", testCase.expectedErrorString, errorString)
					}
				}
			} else {
				if testCase.expectedErrorString != "" {
					t.Errorf("an error was expected, but got none")
				}
			}

		})
	}
}

func TestFromStringHelper(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected simpson
		isValid  bool
	}{
		{
			name:     "valid value homer",
			input:    "Homer Simpson",
			expected: homer,
			isValid:  true,
		},
		{
			name:     "valid value lisa",
			input:    "Lisa Simpson",
			expected: lisa,
			isValid:  true,
		},
		{
			name:     "invalid value",
			input:    "Ned Flanders",
			expected: 0,
			isValid:  false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, ok := cstMeta.FromStringHelper(testCase.input)

			if value != testCase.expected {
				t.Errorf("expected value (int) %d, got %d", testCase.expected, value)
			}

			if ok != testCase.isValid {
				t.Errorf("expected validity %t, got %t", testCase.isValid, ok)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	testCases := []struct {
		name    string
		input   simpson
		isValid bool
	}{
		{
			name:    "valid value homer",
			input:   homer,
			isValid: true,
		},
		{
			name:    "valid value lisa",
			input:   lisa,
			isValid: true,
		},
		{
			name:    "invalid value",
			input:   -1,
			isValid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			isValid := cstMeta.IsValidHelper(testCase.input)

			if isValid != testCase.isValid {
				t.Errorf("expected %t, got %t", testCase.isValid, isValid)
			}
		})
	}
}

func TestMarshalJSONHelper(t *testing.T) {
	// The test will override the metadata, restore them in a proper state
	// after the test.
	jsonStrings := cstMeta.JSONStrings
	t.Cleanup(func() {
		cstMeta.JSONStrings = jsonStrings
	})

	testCases := []struct {
		name         string
		input        simpson
		expected     []byte
		stringSource map[simpson]string
	}{
		{
			name:         "valid value homer",
			input:        homer,
			expected:     []byte(`"Homer Simpson"`),
			stringSource: nil,
		},
		{
			name:         "valid value lisa",
			input:        lisa,
			expected:     []byte(`"Lisa Simpson"`),
			stringSource: nil,
		},
		{
			name:         "invalid value",
			input:        999,
			expected:     nil,
			stringSource: nil,
		},
		{
			name:         "valid value homer (json strings)",
			input:        homer,
			expected:     []byte(`"homer_simpson"`),
			stringSource: jsonStrings,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cstMeta.JSONStrings = testCase.stringSource

			b, err := cstMeta.MarshalJSONHelper(testCase.input)

			if !reflect.DeepEqual(b, testCase.expected) {
				t.Errorf("expected %s, got %s", string(testCase.expected), string(b))
			}

			if testCase.expected != nil && err != nil {
				t.Errorf("unexpected error %v", err)
			}

			if testCase.expected == nil && err == nil {
				t.Errorf("expected an error, got none")
			}
		})
	}
}

func TestUnmarshalJSONHelper(t *testing.T) {
	// The test will override the metadata, restore them in a proper state
	// after the test.
	jsonStrings := cstMeta.JSONStrings
	t.Cleanup(func() {
		cstMeta.JSONStrings = jsonStrings
	})

	testCases := []struct {
		name         string
		input        []byte
		expected     simpson
		stringSource map[simpson]string
	}{
		{
			name:         "valid value homer",
			input:        []byte(`"Homer Simpson"`),
			expected:     homer,
			stringSource: nil,
		},
		{
			name:         "valid value lisa",
			input:        []byte(`"Lisa Simpson"`),
			expected:     lisa,
			stringSource: nil,
		},
		{
			name:         "invalid value",
			input:        []byte(`"Ned Flanders"`),
			expected:     0,
			stringSource: nil,
		},
		{
			name:         "valid value homer (json strings)",
			input:        []byte(`"homer_simpson"`),
			expected:     homer,
			stringSource: jsonStrings,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cstMeta.JSONStrings = testCase.stringSource

			var value simpson
			err := cstMeta.UnmarshalJSONHelper(testCase.input, &value)

			if value != testCase.expected {
				t.Errorf("expected %v, got %v", testCase.expected, value)
			}

			if testCase.expected != 0 && err != nil {
				t.Errorf("unexpected error %v", err)
			}

			if testCase.expected == 0 && err == nil {
				t.Errorf("expected an error, got none")
			}
		})
	}
}
