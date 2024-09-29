package tests_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckScripts(t *testing.T) {
	scriptsDir := "scripts"
	ee, err := os.ReadDir(scriptsDir)
	require.NoError(t, err)
	require.NotNil(t, len(ee))
	for _, e := range ee {
		e := e
		t.Run(e.Name(), func(t *testing.T) {
			scriptName := path.Join(scriptsDir, e.Name())
			t.Log("Script", scriptName)
			require.True(t, e.Type().IsRegular())
			c, err := os.ReadFile(scriptName)
			require.NoError(t, err)
			content := string(c)
			assert.Regexp(t, `^#!/bin/bash\n\n[^\n]`, content)
			assert.NotRegexp(t, `[\t\r\v]`, content) // no tabs etc
			assert.NotRegexp(t, `\x20+\n`, content)  // no leading spaces (except EOF case)
			assert.Regexp(t, `[^\s]\n$`, content)    // strictly one EOL at EOF
		})
	}
}
