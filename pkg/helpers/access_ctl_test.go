package helpers_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/michurin/cnbot/pkg/helpers"
)

type subCase struct {
	check int64
	exp   bool
}

func TestAccessCtl(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name     string
		allow    map[int64]struct{}
		disallow map[int64]struct{}
		expStr   string
		subCase  []subCase
	}{
		{"no_settings", nil, nil, "nobody allowed", []subCase{
			{1, false},
		}},
		{"empty_wl", map[int64]struct{}{}, nil, "white list is empty (nobody can use this bot)", []subCase{
			{1, false},
		}},
		{"wl_1_2", map[int64]struct{}{1: {}, 2: {}}, nil, "1, 2", []subCase{
			{1, true},
			{2, true},
			{3, false},
		}},
		{"empty_bl", nil, map[int64]struct{}{}, "black list is empty (EVERYBODY can use this bot)", []subCase{
			{1, true},
		}},
		{"bl_1", nil, map[int64]struct{}{1: {}}, "all except 1", []subCase{
			{1, false},
			{2, true},
		}},
	} {
		c := c
		a := helpers.NewAccessCtl(c.allow, c.disallow)
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expStr, a.String())
		})
		for _, sc := range c.subCase {
			sc := sc
			t.Run(fmt.Sprintf("%s_%v", c.name, sc.check), func(t *testing.T) {
				assert.Equal(t, sc.exp, a.IsAllowed(sc.check))
			})
		}
	}
}
