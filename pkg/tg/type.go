package tg

import (
	"bytes"
	"errors"
	"strings"
	"unicode/utf8"
)

var (
	// Details: https://en.wikipedia.org/wiki/List_of_file_signatures
	FpGif  = []byte{0x47, 0x49, 0x46, 0x38}
	FpPng  = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	FpJpgA = []byte{0xFF, 0xD8, 0xFF, 0xDB}
	FpJpgB = []byte{0xFF, 0xD8, 0xFF, 0xE0}
	FpJpgC = []byte{0xFF, 0xD8, 0xFF, 0xE1}
	// Errors
	errorMessageTooLong = errors.New("message too long")
	errorInvalidUTF8    = errors.New("invalid UTF8 string")
)

func DataType(data []byte) (
	text string,
	imgExt string,
	err error,
) {
	if bytes.HasPrefix(data, FpGif) {
		imgExt = "gif"
	} else if bytes.HasPrefix(data, FpPng) {
		imgExt = "png"
	} else if bytes.HasPrefix(data, FpJpgA) ||
		bytes.HasPrefix(data, FpJpgB) ||
		bytes.HasPrefix(data, FpJpgC) {
		imgExt = "jpeg"
	} else {
		if utf8.Valid(data) {
			text = strings.TrimSpace(string(data))
			if text == "" {
				text = "(empty)"
			} else if text == "." {
				text = ""
			} else if len(text) > 4000 {
				text = ""
				err = errorMessageTooLong
			}
		} else {
			err = errorInvalidUTF8
		}
	}
	return
}
