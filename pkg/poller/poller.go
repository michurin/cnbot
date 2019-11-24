package poller

import (
	"context"
	"fmt"
	"time"

	"github.com/michurin/cnbot/pkg/api"
	"github.com/michurin/cnbot/pkg/interfaces"
	"github.com/michurin/cnbot/pkg/response"
	"github.com/michurin/cnbot/pkg/workers"
)

func Poller(
	ctx context.Context,
	logger interfaces.Logger,
	a *api.API,
	script string,
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
		result, err := a.JSON(ctx, api.MethodGetUpdates, request)
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
				taskQueue <- workers.Task{
					Text:    u.Message.Text,
					ReplyTo: u.Message.From.ID,
					Script:  script,
				}
			}
			// TODO remove it. It added just for debug
			_, err := a.JSON(ctx, api.MethodSendMessage, map[string]interface{}{
				"chat_id": u.Message.From.ID,
				"text":    "OK",
			})
			if err != nil {
				logger.Log(fmt.Sprintf("Poller error: %s", err))
				sleepWithContext(ctx, 60*time.Second) // TODO sleep flag
				continue
			}
			// TODO /remove it
		}
	}
}

func sleepWithContext(ctx context.Context, d time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}
