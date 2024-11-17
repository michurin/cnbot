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
			// basic types
			"a": nil,
			"b": false,
			"c": true,
			"d": float64(1),
			"e": "text",
			"f": []any{"element"},
			"g": map[string]any{"h": "sub"},
			// complex nested structure
			"i": []any{
				map[string]any{
					"a": float64(1),
					"b": float64(2),
				},
				map[string]any{
					"c": []any{float64(3), float64(4)},
				},
			},
			// corner cases
			"j": "",
			"k": float64(.3),
			"l": float64(.2) + float64(.1),
			// corner cases: partially skipping
			"m": []any{float64(1), nil, float64(3)},                         // m[1] won't appear
			"n": map[string]any{"a": float64(1), "b": nil, "c": float64(3)}, // n["b"] won't appear
			// not appears: empty structures
			"o": map[string]any{},
			"p": []any{},
			// not appears: empty nested structures
			"q": map[string]any{"a": nil, "b": []any{}, "c": map[string]any{}},
			"r": []any{nil, []any{}, map[string]any{}},
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
			"tg_i_1_c=tg_i_1_c_0 tg_i_1_c_1",
			"tg_i_1_c_0=3",
			"tg_i_1_c_1=4",
			"tg_j=",
			"tg_k=0.3",
			"tg_l=0.30000000000000004",
			"tg_m=tg_m_0 tg_m_2",
			"tg_m_0=1",
			"tg_m_2=3",
			"tg_n_a=1",
			"tg_n_c=3",
		}, env)
	})
	t.Run("invalidType", func(t *testing.T) {
		x := float32(1)
		env, err := xjson.JSONToEnv(x)
		assert.EqualError(t, err, "invalid type [pfx=tg]: float32")
		require.Nil(t, env)
	})
	t.Run("invalidTypeInSlice", func(t *testing.T) {
		x := []any{float32(1)}
		env, err := xjson.JSONToEnv(x)
		assert.EqualError(t, err, "invalid type [pfx=tg_0]: float32")
		require.Nil(t, env)
	})
	t.Run("invalidTypeInMap", func(t *testing.T) {
		x := map[string]any{"k": float32(1)}
		env, err := xjson.JSONToEnv(x)
		assert.EqualError(t, err, "invalid type [pfx=tg_k]: float32")
		require.Nil(t, env)
	})
}
