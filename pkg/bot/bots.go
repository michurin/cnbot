package bot

import (
	"context"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

func Bots(ctx context.Context, cfgs []hps.BotConfig) (map[string]hps.BotConfig, error) {
	b := map[string]hps.BotConfig{}
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
		b[botName] = c
		hps.Log(ctx, "Bot configured:", botID, botName)
	}
	return b, nil
}
