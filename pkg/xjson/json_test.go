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
			"m": float64(1<<53 - 1),
			// corner cases: partially skipping
			"n": []any{float64(1), nil, float64(3)},                         // m[1] won't appear
			"o": map[string]any{"a": float64(1), "b": nil, "c": float64(3)}, // n["b"] won't appear
			// not appears: empty structures
			"p": map[string]any{},
			"q": []any{},
			// not appears: empty nested structures
			"r": map[string]any{"a": nil, "b": []any{}, "c": map[string]any{}},
			"s": []any{nil, []any{}, map[string]any{}},
			"t": "line1\nline2",
			"u": "\n\x20line\t\n",
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
			"tg_m=9007199254740991", // 2**53-1 (all bits are '1')
			"tg_n=tg_n_0 tg_n_2",
			"tg_n_0=1",
			"tg_n_2=3",
			"tg_o_a=1",
			"tg_o_c=3",
			"tg_t=line1\nline2",
			"tg_u=\n\x20line\t\n",
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
