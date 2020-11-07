package helpers_test

import (
	"bytes"
	"testing"

	"github.com/michurin/cnbot/pkg/helpers"
	"github.com/stretchr/testify/assert"
)

func TestMessageType(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name string
		b    string
		ig   bool
		txt  string
		md   bool
	}{
		{"empty", "", false, "_empty_", true},
		{"space", " \n ", false, "_empty_", true},
		{"dot", ".", true, "", false},
		{"dot_space", " \n. ", true, "", false},
		{"text", "text", false, "text", false},
		{"text_space", "\n text \n", false, "\n text \n", false},
		{"pre_empty", "%!PRE", false, "_empty \\(pre mode\\)_", true},
		{"pre_empty_ctl", "%!PRE\n\n", false, "_empty \\(pre mode\\)_", true},
		{"pre_space", "%!PRE\n \n", false, "```\n \n\n```", true},
		{"pre_space_same_line", "%!PRE \n \n", false, "```\n \n \n\n```", true},
		{"pre_one_line", "%!PRE ONE", false, "```\n ONE\n```", true},
		{"pre_escape", "%!PRE_ONE", false, "```\n\\_ONE\n```", true},
		{"md_empty", "%!MARKDOWN", false, "_empty \\(markdown mode\\)_", true},
		{"md_empty_ctl", "%!MARKDOWN\n\n", false, "_empty \\(markdown mode\\)_", true},
		{"md_space", "%!MARKDOWN\n \n", false, " \n", true},
		{"md_space_same_line", "%!MARKDOWN \n \n", false, " \n \n", true},
		{"md_one_line", "%!MARKDOWN ONE", false, " ONE", true},
	} {
		t.Run(c.name, func(t *testing.T) {
			ig, txt, md, err := helpers.MessageType([]byte(c.b))
			assert.Nil(t, err)
			assert.Equal(t, c.ig, ig)
			assert.Equal(t, c.txt, txt)
			assert.Equal(t, c.md, md)
		})
	}
	for _, c := range []struct {
		name string
		b    []byte
	}{
		{"err_too_long", bytes.Repeat([]byte{32}, 5000)},
		{"err_utf8", []byte{255}},
	} {
		t.Run(c.name, func(t *testing.T) {
			ig, txt, md, err := helpers.MessageType(c.b)
			assert.NotNil(t, err)
			assert.Equal(t, true, ig)
			assert.Equal(t, "", txt)
			assert.Equal(t, false, md)
		})
	}
}
