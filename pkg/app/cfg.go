package app

import (
	"context"
	"os"

	"github.com/michurin/systemd-env-file/sdenv"

	"github.com/michurin/cnbot/pkg/ctxlog"
	"github.com/michurin/cnbot/pkg/xcfg"
	"github.com/michurin/cnbot/pkg/xlog"
)

func LoadConfigs(files ...string) (map[string]xcfg.Config, string, error) {
	ctx := xlog.Comp(context.Background(), "cfg")
	envs := sdenv.NewCollectsion()
	envs.PushStd(os.Environ())
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, "", ctxlog.Errorfx(ctx, "reading: %w", err)
		}
		pairs, err := sdenv.Parser(data)
		if err != nil {
			return nil, "", ctxlog.Errorfx(ctx, "parser: %w", err)
		}
		envs.Push(pairs)
	}
	cfg, tgAPIOrigin := xcfg.Cfg(ctx, envs.CollectionStd(), files)
	return cfg, tgAPIOrigin, nil
}
