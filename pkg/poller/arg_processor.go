package poller

import (
	"regexp"
	"strings"
)

type ArgProcessor func(text string) ([]string, error)

var reg = regexp.MustCompile("[^a-zA-Z0-9_-]+")

func SafeArgs(text string) ([]string, error) {
	a := strings.Fields(text)
	for i, t := range a {
		a[i] = reg.ReplaceAllString(t, "")
	}
	return a, nil
}
