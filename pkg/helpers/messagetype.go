package helpers

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

var /* const */ markDownEscaping = regexp.MustCompile("([_*[\\]()~`>#+\\-=|{}.!\\\\])")

const markDownMarker = "```"

// It is slightly ugly mix of processor, validator... not just type pure type detector (as ImageType is)
// It has to be rewrote if it grow.
func MessageType(data []byte) (
	ignoreIt bool,
	text string,
	isMarkdown bool,
	err error,
) {
	if !utf8.Valid(data) {
		err = errors.New("invalid message: neither gif/png/jpeg image nor valid UTF8 string")
		return
	}
	text = strings.TrimSpace(string(data))
	if len(text) > 4096 {
		// according documentation this limit applies after entities parsing
		err = errors.New("message too long")
		return
	}
	if text == "" {
		text = "(empty)"
		return
	}
	if text == "." {
		ignoreIt = true
		return
	}
	if len(text) > 6 {
		if strings.HasPrefix(text, markDownMarker) && strings.HasSuffix(text, markDownMarker) {
			isMarkdown = true
			text = strings.TrimSpace(text[3 : len(text)-3])
			text = markDownEscaping.ReplaceAllString(text, "\\$1")
			text = markDownMarker + text + markDownMarker
			return
		}
	}
	return
}
