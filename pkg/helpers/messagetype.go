package helpers

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var /* const */ markDownEscaping = regexp.MustCompile("([_*[\\]()~`>#+\\-=|{}.!\\\\])")

// It is slightly ugly mix of processor, validator... not just type pure type detector (as ImageType is)
// It has to be rewritten if it grow.
//
// Recognize %!PRE, %!MARKDOWN. TODO: %!JSON
func MessageType(data []byte) (
	ignoreIt bool,
	text string,
	isMarkdown bool,
	err error,
) {
	if !utf8.Valid(data) {
		err = errors.New("invalid message: valid UTF8 string")
		return
	}
	text = strings.TrimSpace(string(data))
	if len(text) > 4096 {
		// according documentation this limit applies after entities parsing
		// however this limit is for messages only, for example, image captures has different limitations
		err = errors.New("message too long")
		return
	}
	if text == "" {
		isMarkdown = true
		text = "_empty_"
		return
	}
	if text == "." {
		ignoreIt = true
		return
	}
	if strings.HasPrefix(text, "%!PRE") {
		isMarkdown = true
		text = strings.TrimLeftFunc(text[5:], unicode.IsControl)
		text = markDownEscaping.ReplaceAllString(text, "\\$1")
		if text != "" {
			text = "```\n" + text + "\n```"
		} else {
			text = "_empty \\(pre mode\\)_"
		}
		return
	}
	if strings.HasPrefix(text, "%!MARKDOWN") {
		isMarkdown = true
		text = strings.TrimLeftFunc(text[10:], unicode.IsControl)
		if text == "" {
			text = "_empty \\(markdown mode\\)_"
		}
		return
	}
	return
}
