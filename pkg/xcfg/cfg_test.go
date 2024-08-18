package xcfg_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/michurin/cnbot/pkg/xcfg"
)

func TestCfg(t *testing.T) {
	cfg, tgOrigin := xcfg.Cfg(context.Background(), []string{
		"",
		"x",
		"x=",
		"a=b",
		"a=b=c",

		"tb_x=1", // wrong suffix

		"tb_custom_name_ctrl_addr=c:9090", // correct custom section with file
		"tb_custom_name_token=@testdata/token.txt",
		"tb_custom_name_script=c-short.sh",
		"tb_custom_name_long_running_script=c-worker.sh",

		"tb_noname_ctrl_addr=c:9090", // entire section will be skipped due to file loading error
		"tb_noname_token=@testdata/NOT_EXISTS",
		"tb_noname_script=c-short.sh",
		"tb_noname_long_running_script=c-worker.sh",

		"tb_ctrl_addr=:9090", // correct default section
		"tb_token=xxx",
		"tb_script=short-WILL-BE-OVERRIDDEN.sh", // tb__script
		"tb_long_running_script=worker.sh",

		"tb__script=short.sh", // considering as default, overriding
		"tb_c_script=short.sh",
	}, []string{"/etc/bot.env", "/tmp/bot.env"})
	assert.Equal(t, map[string]xcfg.Config{
		"custom_name": {
			ControlAddr:       "c:9090",
			Token:             "it_is_token_from_file",
			Script:            "c-short.sh",
			LongRunningScript: "c-worker.sh",
			ConfigFileDir:     "/etc", // fake first configuration file's directory
		},
		"default": {
			ControlAddr:       ":9090",
			Token:             "xxx",
			Script:            "short.sh",
			LongRunningScript: "worker.sh",
			ConfigFileDir:     "/etc",
		},
	}, cfg)
	assert.Equal(t, "https://api.telegram.org", tgOrigin)
}
