package tg

import (
	"context"

	hps "github.com/michurin/cnbot/pkg/helpers"
)

type Bot struct {
	Token    string
	Username string
	Script   string
}

func Bots(ctx context.Context, cfgs []hps.BotConfig) ([]Bot, error) {
	b := make([]Bot, len(cfgs))
	for i, c := range cfgs {
		req, err := Encode(c.Token, EncodeGetMe())
		if err != nil {
			hps.Log(ctx, err)
			return nil, err
		}
		out, err := hps.Do(ctx, req)
		if err != nil {
			hps.Log(ctx, err)
			return nil, err
		}
		botID, botName, err := DecodeGetMe(out)
		if err != nil {
			hps.Log(ctx, err)
			return nil, err
		}
		b[i].Token = c.Token
		b[i].Username = botName
		b[i].Script = c.Script
		hps.Log(ctx, "Bot configured:", botID, botName)
	}
	return b, nil
}
