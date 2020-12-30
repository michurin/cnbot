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
		{"pre_empty", "%!PRE", false, "_empty_", true, nil},
		{"pre_empty_ctl", "%!PRE\n\n", false, "_empty_", true, nil},
		{"pre_space", "%!PRE\n \n", false, "_empty_", true, nil},
		{"md", "%!MARKDOWN\ntext", false, "text", true, nil},
		{"md_empty", "%!MARKDOWN", false, "_empty_", true, nil},
		{"md_empty_ctl", "%!MARKDOWN\n\n", false, "_empty_", true, nil},
		{"text_cb", "%!CALLBACK x txt\ntext", false, "text", false, [][][2]string{{{"x", "txt"}}}},
		{"text_cb_pre", "%!CALLBACK y txt2\n%!PRE\ntext", false, "```\ntext\n```", true, [][][2]string{{{"y", "txt2"}}}},
		{"text_cb_pre_nl", "%!CALLBACK z txt3\n\n%!PRE\ntext", false, "```\ntext\n```", true, [][][2]string{{{"z", "txt3"}}}},
		{"text_cb_two", "%!CALLBACK A B\n%!CALLBACK P Q\ntext", false, "text", false, [][][2]string{{{"A", "B"}, {"P", "Q"}}}},
		{"text_cb_two_lines", "%!CALLBACK A B\n%!CALLBACK \n%!CALLBACK P Q\ntext", false, "text", false, [][][2]string{{{"A", "B"}}, {{"P", "Q"}}}},
		{"text_cb_one_word", "%!CALLBACK x\ntext", false, "text", false, [][][2]string{{{"x", "x"}}}},
		{"text_cb_no_message", "%!CALLBACK x txt", false, "_empty_", true, nil},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			ig, txt, md, up, mu, cbt, isa, err := helpers.MessageType([]byte(c.b))
			assert.Nil(t, err)
			assert.False(t, up) // TODO
			assert.Equal(t, c.mu, mu)
			assert.Equal(t, c.ig, ig)
			assert.Equal(t, c.txt, txt)
			assert.Equal(t, c.md, md)
			assert.Empty(t, cbt) // TODO
			assert.Empty(t, isa) // TODO
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
			ig, txt, md, up, mu, cbt, isa, err := helpers.MessageType(c.b)
			assert.NotNil(t, err)
			assert.Nil(t, mu)
			assert.False(t, up)
			assert.Equal(t, true, ig)
			assert.Equal(t, "", txt)
			assert.Equal(t, false, md)
			assert.Empty(t, cbt)
			assert.Empty(t, isa)
		})
	}
}
