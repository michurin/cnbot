package helpers_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/michurin/cnbot/pkg/helpers"
)

type testCase struct {
	Name    string
	Message string
	Exp     struct {
		Text         string
		Ignore       bool
		Markdown     bool
		Update       bool
		Markup       [][][2]string
		CallbackText string `yaml:"callback_text"`
		IsAlert      bool   `yaml:"is_alert"`
	}
}

func TestMessageType_ok(t *testing.T) {
	data, err := ioutil.ReadFile("test_data/message_type.yaml")
	require.NoError(t, err)
	cc := []testCase(nil)
	err = yaml.Unmarshal(data, &cc)
	require.NoError(t, err)
	for _, c := range cc {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			ig, txt, md, up, mu, cbt, isa, err := helpers.MessageType([]byte(c.Message))
			assert.NoError(t, err)
			assert.Equal(t, c.Exp.Text, txt)
			assert.Equal(t, c.Exp.Ignore, ig)
			assert.Equal(t, c.Exp.Markdown, md)
			assert.Equal(t, c.Exp.Update, up)
			assert.Equal(t, c.Exp.Markup, mu)
			assert.Equal(t, c.Exp.CallbackText, cbt)
			assert.Equal(t, c.Exp.IsAlert, isa)
		})
	}
}

func TestMessageType_error(t *testing.T) {
	t.Parallel()
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
