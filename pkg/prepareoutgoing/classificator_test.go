package prepareoutgoing

import (
	"errors"
	"testing"
)

type testCases struct {
	Title     string
	Data      []byte
	IsEmpty   bool
	LeftIt    bool
	RawJSON   bool
	IsImage   bool
	ImageType string
	Error     error
}

func TestClassifier(t *testing.T) {
	for _, c := range []testCases{
		{
			"EmptyInput",
			[]byte{},
			true,
			false,
			false,
			false,
			"",
			nil,
		},
		{
			"Dot:Simplest",
			[]byte("."),
			false,
			true,
			false,
			false,
			"",
			nil,
		},
		{
			"Dot:WithSpaces",
			[]byte(" \t\n\r. \t\n\r"),
			false,
			true,
			false,
			false,
			"",
			nil,
		},
		{
			"Dot:False",
			[]byte(" \t\n\r. \t\n\rXX"),
			false,
			false,
			false,
			false,
			"",
			nil,
		},
		{
			"JSON:Simplest",
			[]byte(`{"one":1}`),
			false,
			false,
			true,
			false,
			"",
			nil,
		},
		{
			"JSON:Deep",
			[]byte(`{"one":{"who":"me"}}`),
			false,
			false,
			true,
			false,
			"",
			nil,
		},
		{
			"JSON:WithSpaces", // newlines!
			[]byte("\n{\n\"one\": {\n\"who\":\"me\" } } "),
			false,
			false,
			true,
			false,
			"",
			nil,
		},
		{
			"Error:InvalidUTF8",
			[]byte{0xff, 0xff},
			false,
			false,
			false,
			false,
			"",
			errors.New("Invalid UTF8 string"),
		},
		{
			"Error:InvalidUTF8:BrokenPNG",
			[]byte{0x89, 0x50},
			false,
			false,
			false,
			false,
			"",
			errors.New("Invalid UTF8 string"),
		},
		{
			"Error:TooLongMessage",
			make([]byte, 10000),
			false,
			false,
			false,
			false,
			"",
			errors.New("Message too long"),
		},
		{
			"Image:GIF",
			[]byte("GIF8"),
			false,
			false,
			false,
			true,
			"gif",
			nil,
		},
	} {
		t.Run(c.Title, func(t *testing.T) {
			isEmpty, leftIt, rawJSON, isImage, imageType, err := classifyData(c.Data)
			assertBool(t, "isEmpty", c.IsEmpty, isEmpty)
			assertBool(t, "leftIt", c.LeftIt, leftIt)
			assertBool(t, "rawJSON", c.RawJSON, rawJSON)
			assertBool(t, "isImage", c.IsImage, isImage)
			assertString(t, "imageType", c.ImageType, imageType)
			assertError(t, c.Error, err)
		})
	}
}

func assertBool(t *testing.T, name string, expected bool, got bool) {
	if expected != got {
		t.Errorf("Invalid %s: expected %t, got %t", name, expected, got)
	}
}

func assertString(t *testing.T, name string, expected string, got string) {
	if expected != got {
		t.Errorf("Invalid %s: expected %s, got %s", name, expected, got)
	}
}

func assertError(t *testing.T, expected error, got error) {
	if got == nil {
		if got != nil {
			t.Errorf("It's have to be error")
		}
	} else {
		if expected == nil {
			t.Errorf("It's have to be NO error")
		} else {
			if got.Error() != expected.Error() {
				t.Errorf(
					"Error strings not matched: expected='%s' but got='%s'",
					expected.Error(),
					got.Error(),
				)
			}
		}
	}
}
