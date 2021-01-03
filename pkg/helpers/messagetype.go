package helpers

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var /* const */ markDownEscaping = regexp.MustCompile("([_*[\\]()~`>#+\\-=|{}.!\\\\])")

var /* const */ labels = []struct {
	label string
	len   int
	re    *regexp.Regexp
}{
	{"markdown", 10, regexp.MustCompile(`^%!MARKDOWN([\r\n]+|$)`)},
	{"pre", 5, regexp.MustCompile(`^%!PRE([\r\n]+|$)`)},
	{"update", 8, regexp.MustCompile(`^%!UPDATE([\r\n]+|$)`)},
	{"callback", 10, regexp.MustCompile(`^%!CALLBACK[^\r\n]*([\r\n]+|$)`)},
	{"callback_text", 6, regexp.MustCompile(`^%!TEXT[^\r\n]+([\r\n]+|$)`)},
	{"callback_alert", 7, regexp.MustCompile(`^%!ALERT[^\r\n]+([\r\n]+|$)`)},
}

func extractLabels(a string) ([][2]string, string) {
	lbs := [][2]string(nil)
	for {
		stop := true
		for _, d := range labels {
			x := d.re.FindString(a)
			l := len(x)
			if l > 0 {
				lbs = append(lbs, [2]string{d.label, strings.TrimSpace(a[d.len:l])})
				a = a[l:]
				stop = false
				break
			}
		}
		if stop {
			break
		}
	}
	return lbs, a
}

func appendNotEmpty(a [][][2]string, b [][2]string) [][][2]string {
	if len(b) > 0 {
		return append(a, b)
	}
	return a
}

func callbackPair(s string) [2]string {
	idx := strings.IndexFunc(s, unicode.IsSpace)
	if idx <= 0 {
		return [2]string{s, s}
	}
	return [2]string{s[:idx], strings.TrimSpace(s[idx:])}
}

// It is slightly ugly mix of processor, validator... not just pure type detector (as ImageType is)
//
// Recognize %!PRE, %!MARKDOWN, %!CALLBACK, %!UPDATE, %!TEXT, %!ALERT
//
// The structure of message is to be:
// - "%!XXX"-labels in any order
// - message
func MessageType(data []byte) (
	ignoreIt bool,
	text string,
	isMarkdown bool,
	forUpdate bool,
	markup [][][2]string,
	callbackText string,
	callbackIsAlert bool,
	err error,
) {
	if !utf8.Valid(data) {
		err = errors.New("invalid message: valid UTF8 string")
		ignoreIt = true
		return
	}
	labels, text := extractLabels(string(data))
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
	originalText := text // before PRE if any
	m := [][2]string(nil)
	for _, l := range labels {
		switch l[0] {
		case "pre":
			isMarkdown = true
			text = "```\n" + markDownEscaping.ReplaceAllString(text, "\\$1") + "\n```"
		case "markdown":
			isMarkdown = true
		case "update":
			forUpdate = true
		case "callback":
			if len(l[1]) == 0 {
				markup = appendNotEmpty(markup, m)
				m = nil
			} else {
				m = append(m, callbackPair(l[1]))
			}
		case "callback_text":
			callbackText = l[1]
			callbackIsAlert = false
		case "callback_alert":
			callbackText = l[1]
			callbackIsAlert = true
		default:
			panic("Unknown label " + l[0])
		}
	}
	markup = appendNotEmpty(markup, m)
	switch strings.TrimSpace(originalText) {
	case "":
		isMarkdown = true
		text = "_empty_"
	case ".":
		ignoreIt = true
		text = ""
	}
	return
}
