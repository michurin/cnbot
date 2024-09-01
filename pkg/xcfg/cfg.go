package xcfg

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/michurin/cnbot/pkg/xlog"
)

var (
	varSfxs = []string{
		"ctrl_addr",
		"token",
		"long_running_script", // order matters; to be greedy longer has to be at beginning
		"script",
	}
	varAllowedSfxs = strings.Join(varSfxs, ", ")
)

const (
	varPrefix    = "tb_"
	varPrefixLen = len(varPrefix)
)

type Config struct {
	ControlAddr       string
	Token             string
	Script            string
	LongRunningScript string
	ConfigFileDir     string
}

func Cfg(ctx context.Context, osEnviron, configFiles []string) (map[string]Config, string) { //nolint:gocognit // slightly inefficient, however it runs only once
	tgAPIOrigin := "https://api.telegram.org"
	x := map[string]map[string]string{}
	for _, pair := range osEnviron {
		ek, ev, ok := strings.Cut(pair, "=")
		if !ok {
			xlog.L(ctx, fmt.Errorf("skipping %q: cannot find `=`", pair))
			continue
		}
		ek = strings.ToLower(ek)
		if len(ev) == 0 {
			xlog.L(ctx, fmt.Errorf("skipping %q: value is empty", pair))
			continue
		}
		if ek == varPrefix+"api_origin" {
			tgAPIOrigin = ev
			continue
		}
		if !strings.HasPrefix(ek, varPrefix) {
			continue
		}
		sfxNotFound := true
		for _, sfx := range varSfxs {
			if strings.HasSuffix(ek, "_"+sfx) {
				sfxNotFound = false
				k := "default"
				if len(ek) > varPrefixLen+len(sfx)+1 {
					k = strings.ToLower(ek[varPrefixLen : len(ek)-1-len(sfx)])
				}
				t := x[k]
				if t == nil {
					t = map[string]string{}
				}
				if x, ok := t[sfx]; ok {
					xlog.L(ctx, fmt.Errorf("overriding %q by %q: %q", x, ev, pair))
				}
				t[sfx] = ev
				x[k] = t
				break
			}
		}
		if sfxNotFound {
			xlog.L(ctx, fmt.Errorf("skipping %q: has TB prefix, but wrong suffix. Allowed: %s", pair, varAllowedSfxs))
		}
	}
	configFileDir := "."
	if len(configFiles) > 0 {
		configFileDir = path.Dir(configFiles[0])
	}
	res := map[string]Config{}
	for k, v := range x {
		if len(v) != 4 {
			xlog.L(ctx, fmt.Errorf("skipping bot name %q: incomplete set of options", k))
			continue
		}
		c := Config{
			ControlAddr:       v[varSfxs[0]],
			Token:             v[varSfxs[1]],
			Script:            v[varSfxs[3]],
			LongRunningScript: v[varSfxs[2]],
			ConfigFileDir:     configFileDir,
		}
		if strings.HasPrefix(c.Token, "@") {
			x, err := os.ReadFile(c.Token[1:])
			if err != nil {
				xlog.L(ctx, fmt.Errorf("skipping bot name %q: cannot get token from file: %q: %w", k, c.Token, err))
				continue
			}
			c.Token = strings.TrimSpace(string(x))
		}
		res[k] = c
	}
	return res, tgAPIOrigin
}
