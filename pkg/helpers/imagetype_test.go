package helpers_test

import (
	"testing"

	"github.com/michurin/cnbot/pkg/helpers"
	"github.com/stretchr/testify/assert"
)

func TestImageType(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name string
		data string
		tp   string
	}{
		{"text", "text", ""},
		{"html", "<B>", ""},
		{"bmp", "BM", ""}, // ignore false triggering
		{"jpeg", "\xFF\xD8\xFF", "jpeg"},
		{"png", "\x89PNG\x0D\x0A\x1A\x0A", "png"},
		{"gif", "GIF89a", "gif"},
	} {
		c := c
		t.Run(c.name, func(t *testing.T) {
			tp := helpers.ImageType([]byte(c.data))
			assert.Equal(t, c.tp, tp)
		})
	}
}
