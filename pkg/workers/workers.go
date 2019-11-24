package workers

import (
	"context"
	"fmt"

	"github.com/michurin/cnbot/pkg/api"
	"github.com/michurin/cnbot/pkg/interfaces"
)

type Task struct {
	Text    string
	ReplyTo int
	Script  string
}

func QueueProcessor(
	ctx context.Context,
	logger interfaces.Logger,
	executor interfaces.Executor,
	taskQueue <-chan Task,
	a *api.API,
) error {
	logger.Log("Queue processor started")
	for {
		select {
		case <-ctx.Done():
			return nil
		case task := <-taskQueue:
			out, err := executor.Run(ctx, task.Script, nil, nil)
			fmt.Printf("OUT: %s\n", string(out))
			fmt.Printf("ERR: %+v\n", err)
			fmt.Printf("TASK: %s\n", task.Text)                                 // TODO use logger
			_, err = a.JSON(ctx, api.MethodSendMessage, map[string]interface{}{ // TODO process body? how?
				"chat_id": task.ReplyTo,
				"text":    string(out),
			})
			if err != nil {
				// TODO sleep?
				logger.Log(err)
			}
		}
	}
}
