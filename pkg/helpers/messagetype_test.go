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
		mu   [][][2]string
	}{
		{"empty", "", false, "_empty_", true, nil},
		{"space", " \n ", false, "_empty_", true, nil},
		{"dot", ".", true, "", false, nil},
		{"dot_space", " \n. ", true, "", false, nil},
		{"text", "text", false, "text", false, nil},
		{"text_space", "\n text \n", false, "\n text \n", false, nil},
		{"pre_empty", "%!PRE", false, "_empty \\(pre mode\\)_", true, nil},
		{"pre_empty_ctl", "%!PRE\n\n", false, "_empty \\(pre mode\\)_", true, nil},
		{"pre_space", "%!PRE\n \n", false, "```\n \n\n```", true, nil},
		{"pre_space_same_line", "%!PRE \n \n", false, "```\n \n \n\n```", true, nil},
		{"pre_one_line", "%!PRE ONE", false, "```\n ONE\n```", true, nil},
		{"pre_escape", "%!PRE_ONE", false, "```\n\\_ONE\n```", true, nil},
		{"md_empty", "%!MARKDOWN", false, "_empty \\(markdown mode\\)_", true, nil},
		{"md_empty_ctl", "%!MARKDOWN\n\n", false, "_empty \\(markdown mode\\)_", true, nil},
		{"md_space", "%!MARKDOWN\n \n", false, " \n", true, nil},
		{"md_space_same_line", "%!MARKDOWN \n \n", false, " \n \n", true, nil},
		{"md_one_line", "%!MARKDOWN ONE", false, " ONE", true, nil},
		{"text_cb", "%!CALLBACK x txt\ntext", false, "text", false, [][][2]string{{{"x", "txt"}}}},
		{"text_cb_pre", "%!CALLBACK y txt2\n%!PRE\ntext", false, "```\ntext\n```", true, [][][2]string{{{"y", "txt2"}}}},
		{"text_cb_pre_nl", "%!CALLBACK z txt3\n\n%!PRE\ntext", false, "```\ntext\n```", true, [][][2]string{{{"z", "txt3"}}}},
		{"text_cb_two", "%!CALLBACK A B\n%!CALLBACK P Q\ntext", false, "text", false, [][][2]string{{{"A", "B"}, {"P", "Q"}}}},
		{"text_cb_two_lines", "%!CALLBACK A B\n%!CALLBACK \n%!CALLBACK P Q\ntext", false, "text", false, [][][2]string{{{"A", "B"}}, {{"P", "Q"}}}},
		{"text_cb_one_word", "%!CALLBACK x\ntext", false, "text", false, [][][2]string{{{"x", "x"}}}},
		{"text_cb_no_message", "%!CALLBACK x txt", false, "_empty \\(callback mode\\)_", true, [][][2]string{{{"x", "txt"}}}},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			ig, txt, md, up, mu, err := helpers.MessageType([]byte(c.b))
			assert.Nil(t, err)
			assert.False(t, up) // TODO
			assert.Equal(t, c.mu, mu)
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
		c := c
		t.Run(c.name, func(t *testing.T) {
			ig, txt, md, up, mu, err := helpers.MessageType(c.b)
			assert.NotNil(t, err)
			assert.Nil(t, mu)
			assert.False(t, up)
			assert.Equal(t, true, ig)
			assert.Equal(t, "", txt)
			assert.Equal(t, false, md)
		})
	}
}
