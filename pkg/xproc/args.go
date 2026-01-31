package xproc

import (
	"strings"
	"unicode/utf8"
)

func SanitizeArgs(a []string) []string {
	n := min(len(a), 32) // limit number of arguments
	b := make([]string, n)
	for i := range n {
		if !utf8.ValidString(a[i]) { // skip invalid strings
			continue
		}
		r := []rune(a[i])
		p := make([]rune, 0, len(r))
		for _, c := range r {
			if c <= 0x20 || c >= 127 {
				continue
			}
			switch c {
			case '*', '?', '[', ']', // systemd: GLOB_CHARS
				'"', '\\', '`', '$', // systemd: SHELL_NEED_ESCAPE
				'\'', '(', ')', '<', '>', '|', '&', ';', '!', // systemd: SHELL_NEED_QUOTES
				'/': // extra protection to disable relative and paths, however we keep '.' for IP addresses
				continue
			}
			p = append(p, c)
			if len(p) >= 256 { // one argument limit
				break
			}
		}
		b[i] = strings.ToLower(string(p))
	}
	return b
}
