package bot

import (
	"context"
	"time"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

const pollingRequestTimeOutSeconds = 30
const errorSleepDuration = time.Second * 10

func Poller(baseCtx context.Context, botName string, bot hps.BotConfig, msgQueue chan<- tg.Message) {
	var offset int
	var mm []tg.Message
	hps.Log(hps.Label(baseCtx, botName), "Poller runs for bot", botName)
MainLoop:
	for {
		ctx := hps.Label(baseCtx, hps.RandLabel(), botName)
		select {
		case <-ctx.Done():
			hps.Log(ctx, "Poller is halted by context canceling")
			break MainLoop
		default:
		}
		req, err := tg.EncodeGetUpdates(offset, pollingRequestTimeOutSeconds)
		if err != nil { // in fact, it is reason for panic
			hps.Log(ctx, err)
			hps.Sleep(ctx, errorSleepDuration)
			continue
		}
		out, err := hps.Do(ctx, tg.Encode(bot.Token, req))
		if err != nil {
			hps.Log(ctx, err)
			hps.Sleep(ctx, errorSleepDuration)
			continue
		}
		mm, offset, err = tg.DecodeGetUpdate(out, offset, botName)
		if err != nil {
			hps.Log(ctx, err)
			hps.Sleep(ctx, errorSleepDuration)
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
	hps.Log(baseCtx, "Poller is stopped")
}
