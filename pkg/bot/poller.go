package bot

import (
	"context"
	"time"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

const (
	pollingRequestTimeOutSeconds = 30
	errorSleepDuration           = time.Second * 10
)

func Poller(baseCtx context.Context, botName string, bot hps.BotConfig, msgQueue chan<- tg.Message) {
	var offset int64
	var mm []tg.Message
	var sleep bool
	hps.Log(hps.Label(baseCtx, botName), "Poller runs")
	ctx := baseCtx
MainLoop:
	for {
		if sleep { // comes from previous iteration; it will never execute in first run
			err := hps.Sleep(ctx, errorSleepDuration)
			if err != nil {
				hps.Log(ctx, "Poller is halted by ctx")
				break
			}
			sleep = false
		}
		ctx = hps.Label(baseCtx, hps.RandLabel(), botName)
		req, err := tg.EncodeGetUpdates(offset, pollingRequestTimeOutSeconds)
		if err != nil { // in fact, it is reason for panic
			hps.Log(ctx, err)
			sleep = true
			continue
		}
		out, err := hps.Do(ctx, tg.Encode(bot.Token, req))
		if err != nil {
			hps.Log(ctx, err)
			sleep = true
			continue
		}
		mm, offset, err = tg.DecodeGetUpdates(out, offset, botName)
		if err != nil {
			hps.Log(ctx, err)
			sleep = true
			continue
		}
		for _, m := range mm {
			select {
			case <-ctx.Done():
				hps.Log(ctx, "Poller is halted by context canceling")
				break MainLoop
			case msgQueue <- m:
			}
		}
	}
	hps.Log(hps.Label(baseCtx, botName), "Poller is stopped")
}
