package xjson_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/michurin/cnbot/pkg/xjson"
)

func TestJSONToEnv(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		x := map[string]any{
			"a": nil,
			"b": false,
			"c": true,
			"d": float64(1),
			"e": "text",
			"f": []any{"element"},
			"g": map[string]any{
				"h": "sub",
			},
			"i": []any{
				map[string]any{
					"a": float64(1),
					"b": float64(2),
				},
				map[string]any{
					"c": float64(3),
					"d": float64(4),
				},
			},
		}
		env, err := xjson.JSONToEnv(x)
		require.NoError(t, err)
		assert.Equal(t, []string{
			"tg_b=false",
			"tg_c=true",
			"tg_d=1",
			"tg_e=text",
			"tg_f=tg_f_0",
			"tg_f_0=element",
			"tg_g_h=sub",
			"tg_i=tg_i_0 tg_i_1",
			"tg_i_0_a=1",
			"tg_i_0_b=2",
			"tg_i_1_c=3",
			"tg_i_1_d=4",
		}, env)
	})
	t.Run("invalidType", func(t *testing.T) {
		x := float32(1)
		env, err := xjson.JSONToEnv(x)
		assert.Error(t, err)
		require.Nil(t, env)
	})
	t.Run("invalidTypeInSlice", func(t *testing.T) {
		x := []any{float32(1)}
		env, err := xjson.JSONToEnv(x)
		assert.Error(t, err)
		require.Nil(t, env)
	})
	t.Run("invalidTypeInMap", func(t *testing.T) {
		x := map[string]any{"k": float32(1)}
		env, err := xjson.JSONToEnv(x)
		assert.Error(t, err)
		require.Nil(t, env)
	})
}
