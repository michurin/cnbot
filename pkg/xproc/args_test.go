package xproc_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/michurin/cnbot/pkg/xproc"
)

func TestSanitizeArgs(t *testing.T) {
	t.Run("white_list", func(t *testing.T) {
		x := make([]rune, 1024)
		for i := range 1024 {
			x[i] = rune(i)
		}
		assert.Equal(t,
			[]string{"#%+,-./0123456789:=@abcdefghijklmnopqrstuvwxyz^_abcdefghijklmnopqrstuvwxyz{}~"},
			xproc.SanitizeArgs([]string{string(x)}))
	})
	t.Run("strings", func(t *testing.T) {
		assert.Equal(t,
			[]string{
				"x",
				"",
				"",
				"",
				"_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_123456789_12345",
			},
			xproc.SanitizeArgs([]string{
				"x",                              // simple
				"",                               // empty
				">",                              // empty after cleaning
				string([]byte{48, 255, 255}),     // invalid UTF8
				strings.Repeat("_123456789", 30), // too long
			}))
	})
	t.Run("limit", func(t *testing.T) {
		x := make([]string, 33)
		assert.Len(t, xproc.SanitizeArgs(x), 32)
	})
	t.Run("nil", func(t *testing.T) {
		assert.Equal(t, []string{}, xproc.SanitizeArgs([]string(nil)))
	})
	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, []string{}, xproc.SanitizeArgs([]string{}))
	})
}
