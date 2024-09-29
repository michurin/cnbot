package app

import (
	"context"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/michurin/systemd-env-file/sdenv"

	"github.com/michurin/cnbot/pkg/ctxlog"
	"github.com/michurin/cnbot/pkg/xlog"
)

const (
	varPrefix    = "tb_"
	varPrefixLen = len(varPrefix)
)

var varSfxs = []string{
	"_ctrl_addr",
	"_token",
	"_config_dir",
	"_long_running_script", // order matters; to be greedy longer has to be first
	"_script",
}

func LoadConfigs(files ...string) (map[string]BotConfig, string, error) {
	ctx := xlog.Comp(context.Background(), "cfg")

	// init collection of variables by system environment

	envs := sdenv.NewCollectsion()
	envs.PushStd(os.Environ()) // so environment variables has maximum priority

	// enrich variables from files

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, "", ctxlog.Errorfx(ctx, "reading: %s: %w", file, err)
		}
		pairs, err := sdenv.Parser(data)
		if err != nil {
			return nil, "", ctxlog.Errorfx(ctx, "parser: %s: %w", file, err)
		}
		envs.Push(pairs)
		envs.Push(configFileDirForEachSetup(file, pairs)) // add configuration file information with minimum priority
	}

	// sort to be reproducible

	ev := envs.CollectionStd()
	sort.Strings(ev) // we have to sort in original case, before tolowering

	// consider all

	c, err := buildConfigs(ev)
	if err != nil {
		return nil, "", ctxlog.Errorfx(ctx, "building configs: %w", err)
	}

	return c, buildOrigin(ev), nil
}

func buildOrigin(pairs []string) string {
	for _, pair := range pairs {
		ek, ev, ok := strings.Cut(pair, "=")
		if !ok {
			continue
		}
		if len(ek) <= varPrefixLen { // skip: it cannot be configuration variable
			continue
		}
		if strings.ToLower(ek[varPrefixLen:]) == "api_origin" {
			return ev
		}
	}
	return "https://api.telegram.org"
}

func buildConfigs(pairs []string) (map[string]BotConfig, error) {
	configMaps := map[string][5]string{} // TODO: 5 hardcoded
	// TODO minor code duplication in this loop
	for _, pair := range pairs {
		ek, ev, ok := strings.Cut(pair, "=")
		if !ok {
			continue
		}
		if len(ek) <= varPrefixLen { // skip: it cannot be configuration variable
			continue
		}
		ekLower := strings.ToLower(ek)
		if !strings.HasPrefix(ekLower, varPrefix) { // skip: is not a configuration variable due to prefix
			continue
		}
		for i, sfx := range varSfxs {
			if strings.HasSuffix(ekLower, sfx) {
				k := "default"
				if len(ek) > varPrefixLen+len(sfx) {
					k = ekLower[varPrefixLen : len(ek)-len(sfx)]
				}
				v := configMaps[k]
				v[i] = ev
				configMaps[k] = v
				break
			}
		}
	}
	res := map[string]BotConfig{}
	for k, v := range configMaps {
		token := v[1]
		if token == "" {
			return nil, fmt.Errorf("bot %[1]s: token is empty, you must specify tb_%[1]s_token", k)
		}
		if strings.HasPrefix(token, "@") {
			x, err := os.ReadFile(token[1:]) // TODO consider relative paths; it has to share logic with runner
			if err != nil {
				xlog.L(context.TODO(), fmt.Errorf("skipping bot name %q: cannot get token from file: %q: %w", k, token, err))
				continue
			}
			token = strings.TrimSpace(string(x))
		}
		res[k] = BotConfig{
			ControlAddr:       v[0],
			Token:             token,
			Script:            v[4],
			LongRunningScript: v[3],
			ConfigFileDir:     v[2],
		}
	}
	return res, nil
}

func configFileDirForEachSetup(cfgFile string, pairs [][2]string) [][2]string {
	dir := path.Dir(cfgFile)
	res := [][2]string(nil)
	for _, pair := range pairs {
		ek, ev := pair[0], pair[1]
		if len(ev) <= varPrefixLen { // skip: it cannot be configuration variable
			continue
		}
		ekLower := strings.ToLower(ek)
		if strings.HasPrefix(ekLower, varPrefix) { // skip: is not a configuration variable due to prefix
			continue
		}
		for _, sfx := range varSfxs {
			if strings.HasSuffix(ekLower, sfx) {
				k := "default"
				if len(ek) > varPrefixLen+len(sfx) {
					k = ekLower[varPrefixLen : len(ek)-len(sfx)]
				}
				res = append(res, [2]string{varPrefix + k + "_config_dir", dir})
				break // end on first matching
			}
		}
	}
	return res
}
