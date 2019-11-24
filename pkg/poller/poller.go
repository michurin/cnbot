package poller

import (
	"context"
	"fmt"
	"time"

	"github.com/michurin/cnbot/pkg/api"
	"github.com/michurin/cnbot/pkg/interfaces"
	"github.com/michurin/cnbot/pkg/processors"
	"github.com/michurin/cnbot/pkg/response"
	"github.com/michurin/cnbot/pkg/workers"
)

func Poller(
	ctx context.Context,
	logger interfaces.Logger,
	a *api.API,
	script string,
	messageProcessor processors.MessageProcessor,
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
		result, err := a.Call(ctx, api.MethodGetUpdates, body)
		if err != nil {
			logger.Log(fmt.Sprintf("Poller error: %s", err))
			sleepWithContext(ctx, 60*time.Second) // TODO sleep flag
			continue
		}
		updates, err := response.ParseUpdates(result)
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
				taskQueue <- workers.Task{
					Text:    u.Message.Text,
					Args:    args,
					ReplyTo: u.Message.From.ID,
					Script:  script,
				}
			}
			/* // TODO remove it. It added just for debug
			xbody, _ := api.EncodeJSON(map[string]interface{}{
				"chat_id": u.Message.From.ID,
				"text":    "OK",
			})
			_, err := a.Call(ctx, api.MethodSendMessage, xbody)
			if err != nil {
				logger.Log(fmt.Sprintf("Poller error: %s", err))
				sleepWithContext(ctx, 60*time.Second) // TODO sleep flag
				continue
			}
			*/ // TODO /remove it
		}
	}
}

func sleepWithContext(ctx context.Context, d time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}
