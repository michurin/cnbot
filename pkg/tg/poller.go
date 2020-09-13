package tg

import (
	"context"
	"time"

	hps "github.com/michurin/cnbot/pkg/helpers"
)

const pollingRequestTimeOutSeconds = 10

const errorSleepDuration = time.Second * 10

func errorSleep(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(duration):
	}
}

func Poller(ctx context.Context, bot Bot, msgQueue chan<- Message) {
	var offset int
	var mm []Message
MainLoop:
	for {
		select {
		case <-ctx.Done():
			hps.Log(ctx, "Poller is halted by context canceling")
			break MainLoop
		default:
		}
		req, err := Encode(bot.Token, EncodeGetUpdates(offset, pollingRequestTimeOutSeconds))
		if err != nil {
			hps.Log(ctx, err)
			errorSleep(ctx, errorSleepDuration)
			continue
		}
		out, err := hps.Do(ctx, req)
		if err != nil {
			hps.Log(ctx, err)
			errorSleep(ctx, errorSleepDuration)
			continue
		}
		mm, offset, err = DecodeGetUpdate(out, offset, bot.Username)
		if err != nil {
			hps.Log(ctx, err)
			errorSleep(ctx, errorSleepDuration)
			continue
		}
		for _, m := range mm {
			select {
			case <-ctx.Done():
				hps.Log(ctx, "Poller is halted by context canceling")
				break MainLoop
			case msgQueue <- m:
			}
			hps.Log(ctx, m)
		}
	}
}
