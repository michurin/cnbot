package helpers

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var /* const */ markDownEscaping = regexp.MustCompile("([_*[\\]()~`>#+\\-=|{}.!\\\\])")

// It is slightly ugly mix of processor, validator... not just pure type detector (as ImageType is)
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
		ignoreIt = true
		return
	}
	text = string(data)
	if len(text) > 4096 {
		// TODO ugly check
		// - according documentation this limit applies after entities parsing
		// - this limit is for messages only, for example, image captures has another limitations
		// to perform this check correctly, we have to parse markdown locally; what we don't do yet
		ignoreIt = true
		text = ""
		err = errors.New("message too long")
		return
	}
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		isMarkdown = true
		text = "_empty_"
		return
	}
	if trimmed == "." {
		ignoreIt = true
		text = ""
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
