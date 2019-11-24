package processors

import (
	"regexp"
	"strings"
)

var reg = regexp.MustCompile("[^a-zA-Z0-9_-]+")

func Safe(text string) ([]string, error) {
	a := strings.Fields(text)
	for i, t := range a {
		a[i] = reg.ReplaceAllString(t, "")
	}
	return a, nil
}
