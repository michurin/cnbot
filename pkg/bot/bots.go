package bot

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/michurin/cnbot/pkg/tg"
	hps "github.com/michurin/cnbot/pkg/helpers"
)

type Bot struct {
	Token        string
	AllowedUsers map[int]struct{}
	Script       string
	WorkingDir   string
}

func Bots(ctx context.Context, cfgs []hps.BotConfig) (map[string]Bot, error) {
	b := map[string]Bot{}
	for _, c := range cfgs {
		out, err := hps.Do(ctx, tg.Encode(c.Token, tg.EncodeGetMe()))
		if err != nil {
			hps.Log(ctx, err)
			return nil, err
		}
		botID, botName, err := tg.DecodeGetMe(out)
		if err != nil {
			hps.Log(ctx, err)
			return nil, err
		}
		allowedUsers := map[int]struct{}{}
		for _, u := range c.AllowedUsers {
			if _, ok := allowedUsers[u]; ok {
				hps.Log(ctx, u, errors.New("user is already allowed"))
			}
			allowedUsers[u] = struct{}{}
		}
		workingDir := c.WorkingDir
		if workingDir == "" {
			workingDir = string(filepath.Separator)
			hps.Log(ctx, errors.New("working dir not specified"))
		}
		b[botName] = Bot{
			Token:        c.Token,
			Script:       c.Script,
			AllowedUsers: allowedUsers,
			WorkingDir:   workingDir,
		}
		hps.Log(ctx, "Bot configured:", botID, botName)
	}
	return b, nil
}
