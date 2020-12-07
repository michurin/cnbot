package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractLabels(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name string
		in   string
		lbs  [][2]string
		body string
	}{
		{"simples", "body", nil, "body"},
		{"simples_spaces", "\n body \n", nil, "\n body \n"},
		{"wrong_label", "%!ONE\nbody", nil, "%!ONE\nbody"},
		{"pre_no_newline", "%!PRE", [][2]string{{"pre", ""}}, ""},
		{"pre", "%!PRE\nbody", [][2]string{{"pre", ""}}, "body"},
		{"pre_markdown", "%!PRE\n%!MARKDOWN\nbody", [][2]string{{"pre", ""}, {"markdown", ""}}, "body"},
		{"pre_markdown_spaces", "%!PRE\n\n%!MARKDOWN\n\nbody", [][2]string{{"pre", ""}, {"markdown", ""}}, "body"},
		{"pre_callback", "%!PRE\n%!CALLBACK A B C\nbody", [][2]string{{"pre", ""}, {"callback", "A B C"}}, "body"},
		{"update", "%!UPDATE\nbody", [][2]string{{"update", ""}}, "body"},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			lbs, body := extractLabels(c.in)
			assert.Equal(t, c.lbs, lbs)
			assert.Equal(t, c.body, body)
		})
	}
}
