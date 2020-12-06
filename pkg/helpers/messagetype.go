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
// Recognize %!PRE, %!MARKDOWN, %!CALLBACK, %!UPDATE
//
// The structure of message is to be:
// - Optional %!UPDATE
// - Zero or more %!CALLBACK lines
// - Optional %!PRE or %!MARKDOWN
// - message
func MessageType(data []byte) (
	ignoreIt bool,
	text string,
	isMarkdown bool,
	forUpdate bool,
	markup [][][2]string,
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
	if strings.HasPrefix(text, "%!UPDATE") {
		forUpdate = true
		idx := strings.IndexFunc(text, unicode.IsControl)
		text = strings.TrimLeftFunc(text[idx:], unicode.IsControl)
	}
	m := [][2]string(nil)
	for {
		if strings.HasPrefix(text, "%!CALLBACK") {
			var a string
			idx := strings.IndexFunc(text, unicode.IsControl)
			if idx > 0 {
				a = text[:idx]
				text = strings.TrimLeftFunc(text[idx:], unicode.IsControl)
			} else {
				a = text
				text = ""
			}
			a = a[10:] // remove %!CALLBACK
			a = strings.TrimSpace(a)
			if len(a) == 0 {
				if len(m) > 0 {
					markup = append(markup, m)
				}
				m = [][2]string(nil)
			} else {
				idx := strings.IndexFunc(a, unicode.IsSpace)
				if idx <= 0 {
					m = append(m, [2]string{a, a})
				} else {
					m = append(m, [2]string{a[:idx], strings.TrimSpace(a[idx:])})
				}
			}
		} else {
			break
		}
	}
	if len(m) > 0 {
		markup = append(markup, m)
	}
	if text == "" {
		isMarkdown = true
		text = "_empty \\(callback mode\\)_"
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
