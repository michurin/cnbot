package prepareoutgoing

import (
	"bytes"
	"errors"
	"testing"
)

type testCases struct {
	Title      string
	Data       []byte
	IsEmpty    bool
	LeftIt     bool
	IsRaw      bool
	RawMethod  string
	RawPayload []byte
	IsImage    bool
	ImageType  string
	Error      error
}

func TestClassifier(t *testing.T) {
	for _, c := range []testCases{
		{
			"EmptyInput",
			[]byte{},
			true,
			false,
			false,
			"",
			nil,
			false,
			"",
			nil,
		},
		{
			"EmptyInputBlankChars",
			[]byte{"\r\n\t\x20"},
			true,
			false,
			false,
			"",
			nil,
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
			"",
			nil,
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
			"",
			nil,
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
			"",
			nil,
			false,
			"",
			nil,
		},
		{
			"JSON:Simplest",
			[]byte(`sendMessage{"one":1}`),
			false,
			false,
			true,
			"sendMessage",
			[]byte(`{"one":1}`),
			false,
			"",
			nil,
		},
		{
			"JSON:Deep",
			[]byte(`sendMessage{"one":{"who":"me"}}`),
			false,
			false,
			true,
			"sendMessage",
			[]byte(`{"one":{"who":"me"}}`),
			false,
			"",
			nil,
		},
		{
			"JSON:WithSpaces", // newlines!
			[]byte("\nsendMessage\n{\n\"one\": {\n\"who\":\"me\" } }\n"),
			false,
			false,
			true,
			"sendMessage",
			[]byte("{\n\"one\": {\n\"who\":\"me\" } }"),
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
			"",
			nil,
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
			"",
			nil,
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
			"",
			nil,
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
			"",
			nil,
			true,
			"gif",
			nil,
		},
	} {
		t.Run(c.Title, func(t *testing.T) {
			isEmpty, leftIt, isRaw, rawMethod, rawPayload, isImage, imageType, err := classifyData(c.Data)
			assertBool(t, "isEmpty", c.IsEmpty, isEmpty)
			assertBool(t, "leftIt", c.LeftIt, leftIt)
			assertBool(t, "isRaw", c.IsRaw, isRaw)
			assertString(t, "rawMethod", c.RawMethod, rawMethod)
			assertByte(t, "rawPayload", c.RawPayload, rawPayload)
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

func assertByte(t *testing.T, name string, expected []byte, got []byte) {
	if !bytes.Equal(expected, got) {
		t.Errorf("Invalid %s: expected %s, got %s", name, expected, got)
	}
}

func assertError(t *testing.T, expected error, got error) {
	if got == nil {
		if got != nil {
			t.Error("It's have to be error")
		}
	} else {
		if expected == nil {
			t.Errorf("It's have to be NO error, got: %s", got.Error())
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
