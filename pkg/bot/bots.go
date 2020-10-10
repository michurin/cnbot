package bot

import (
	"context"
	"fmt"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

func Bots(rootCtx context.Context, cfgs map[string]hps.BotConfig) error {
	for nick, c := range cfgs {
		ctx := hps.Label(rootCtx, nick)
		out, err := hps.Do(ctx, tg.Encode(c.Token, tg.EncodeGetMe()))
		if err != nil {
			hps.Log(ctx, err)
			return err
		}
		botID, botName, err := tg.DecodeGetMe(out)
		if err != nil {
			hps.Log(ctx, err)
			return err
		}
		hps.Log(ctx, fmt.Sprintf("Bot %q is OK. ID: %d, Name: %q", nick, botID, botName))
	}
	return nil
}
