package bot

import (
	"context"
	"time"

	"github.com/michurin/minlog"

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
	baseCtx = minlog.Label(baseCtx, botName)
	minlog.Log(baseCtx, "Poller runs")
MainLoop:
	for {
		if sleep { // comes from previous iteration; it will never execute in first run
			err := hps.Sleep(baseCtx, errorSleepDuration)
			if err != nil {
				minlog.Log(baseCtx, "Poller is halted by ctx")
				break
			}
			sleep = false
		}
		ctx := minlog.Label(baseCtx, hps.AutoLabel())
		req, err := tg.EncodeGetUpdates(offset, pollingRequestTimeOutSeconds)
		if err != nil { // in fact, it is reason for panic
			minlog.Log(ctx, err)
			sleep = true
			continue
		}
		out, err := hps.Do(ctx, tg.Encode(bot.Token, req))
		if err != nil {
			minlog.Log(ctx, err)
			sleep = true
			continue
		}
		mm, offset, err = tg.DecodeGetUpdates(out, offset, botName)
		if err != nil {
			minlog.Log(ctx, err)
			sleep = true
			continue
		}
		for _, m := range mm {
			select {
			case <-ctx.Done():
				minlog.Log(ctx, "Poller is halted by context canceling")
				break MainLoop
			case msgQueue <- m:
			}
		}
	}
	minlog.Log(baseCtx, "Poller is stopped")
}
