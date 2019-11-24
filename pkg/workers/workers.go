package workers

import (
	"context"
	"fmt"
	"time"

	"github.com/michurin/cnbot/pkg/api"
	"github.com/michurin/cnbot/pkg/execute"
	"github.com/michurin/cnbot/pkg/interfaces"
)

type Task struct {
	Text    string // TODO what is this field fore?
	Args    []string
	ReplyTo int
	Script  string
}

func QueueProcessor(
	ctx context.Context,
	logger interfaces.Logger,
	runner *execute.Executor,
	taskQueue <-chan Task,
	a *api.API,
) error {
	logger.Log("Queue processor started")
	for {
		select {
		case <-ctx.Done():
			return nil
		case task := <-taskQueue:
			out, err := runner.Run(ctx, execute.ScriptInfo{
				Name:    task.Script,
				Timeout: 2 * time.Second, // TODO get from Task
				Env:     nil,             // TODO
				Args:    task.Args,
			})
			fmt.Printf("OUT: %s\n", string(out))
			fmt.Printf("ERR: %+v\n", err)
			fmt.Printf("TASK: %s\n", task.Text) // TODO use logger

			/* How to send photo
			body, err := api.EncodeMultipart(task.ReplyTo, out)
			if err != nil {
				panic(err)
			}
			_, err = a.Call(ctx, api.MethodSendPhoto, body)
			if err != nil {
				panic(err)
			}
			continue
			*/

			body, err := api.EncodeJSON(map[string]interface{}{
				"chat_id": task.ReplyTo,
				"text":    string(out),
			})
			if err != nil {
				// TODO sleep?
				logger.Log(err)
			}
			_, err = a.Call(ctx, api.MethodSendMessage, body)
			if err != nil {
				// TODO sleep?
				logger.Log(err)
			}
		}
	}
}
