package poller

import (
	"context"
	"fmt"
	"time"

	"github.com/michurin/cnbot/pkg/api"
	"github.com/michurin/cnbot/pkg/interfaces"
	"github.com/michurin/cnbot/pkg/workers"
)

func Poller(
	ctx context.Context,
	logger interfaces.Logger,
	botName string,
	messageFilter MessageFilter,
	pollAPI api.Interface,
	replierAPI api.Interface,
	env []string,
	script string,
	messageProcessor ArgProcessor,
	taskQueue chan<- workers.Task,
) error {
	logger.Log("Poller started")
	offset := 0
	addOffset := false
	for {
		select {
		case <-ctx.Done():
			logger.Log("Break!")
			return nil
		default:
		}
		request := map[string]interface{}{
			"timeout":         55, // TODO param
			"allowed_updates": []string{"message"},
		}
		if addOffset {
			request["offset"] = offset + 1
		}
		body, err := api.EncodeJSON(request)
		if err != nil {
			logger.Log(fmt.Sprintf("Poller error: %s", err))
			sleepWithContext(ctx, 60*time.Second) // TODO sleep flag
			continue
		}
		result, err := pollAPI.Call(ctx, api.MethodGetUpdates, body)
		if err != nil {
			logger.Log(fmt.Sprintf("Poller error: %s", err))
			sleepWithContext(ctx, 60*time.Second) // TODO sleep flag
			continue
		}
		updates, err := ParseUpdates(result)
		if err != nil {
			panic("TODO") // TODO ENORMOUS TODO
		}

		for _, u := range updates {
			itIsNewMessage := false // In fact, Telegram could send same message twice
			if addOffset {
				if offset < u.UpdateID {
					offset = u.UpdateID
					itIsNewMessage = true
				}
			} else {
				offset = u.UpdateID
				addOffset = true
				itIsNewMessage = true
			}
			if itIsNewMessage {
				args, xerr := messageProcessor(u.Message.Text)
				if xerr != nil {
					logger.Log(fmt.Sprintf("Poller error: %s", err))
					continue // TODO what we have to do here?
				}
				if messageFilter.IsAllowed(u) {
					taskQueue <- workers.Task{
						Args:        args,
						BotName:     botName,
						ReplyTo:     u.Message.From.ID,
						ReplayToAPI: replierAPI,
						Env:         env,
						Script:      script,
					}
				} else {
					logger.Log(fmt.Sprintf("User %d DISALLOWED!", u.Message.From.ID))
				}
			}
		}
	}
}

func sleepWithContext(ctx context.Context, d time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}
