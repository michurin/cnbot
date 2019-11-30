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
	BotName string
	ReplyTo int
	Script  string
}

func QueueProcessor(
	ctx context.Context,
	logger interfaces.Logger,
	runner *execute.Executor,
	taskQueue <-chan Task,
	a map[string]*api.API,
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
				Env: []string{
					fmt.Sprintf("T_USER=%d", task.ReplyTo),
					fmt.Sprintf("T_BOT=%s", task.BotName),
				},
				Args: task.Args,
			})
			if err != nil {
				// TODO sleep? break?
				logger.Log(err)
			}
			fmt.Printf("OUT: %q\n", string(out))
			fmt.Printf("ERR: %+v\n", err)
			fmt.Printf("TASK: %s\n", task.Text) // TODO use logger

			method, body, err := api.SimpleRequest(out, task.ReplyTo)
			if err != nil {
				// TODO sleep? break?
				logger.Log(err)
			} else {
				// TODO refactor it
				_, err = a[task.BotName].Call(ctx, method, body) // TODO check botName exists?
				if err != nil {
					// TODO sleep?
					logger.Log(err)
				}
			}
		}
	}
}
