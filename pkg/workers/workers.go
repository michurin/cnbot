package workers

import (
	"context"
	"fmt"
	"time"

	"github.com/michurin/cnbot/pkg/apirequest"
	"github.com/michurin/cnbot/pkg/execute"
	"github.com/michurin/cnbot/pkg/interfaces"
)

type Task struct {
	Args        []string
	BotName     string
	ReplyTo     int
	Env         []string
	Script      string
	ReplayToAPI interfaces.Interface
}

func QueueProcessor(
	ctx context.Context,
	logger interfaces.Logger,
	runner *execute.Executor,
	taskQueue <-chan Task,
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
				Env: append(
					task.Env,
					fmt.Sprintf("T_USER=%d", task.ReplyTo),
					fmt.Sprintf("T_BOT=%s", task.BotName)),
				Args: task.Args,
			})
			if err != nil {
				// TODO sleep? break?
				logger.Log(err)
			}
			fmt.Printf("OUT: %q\n", string(out))
			fmt.Printf("ERR: %+v\n", err)

			method, body, err := apirequest.SimpleRequest(out, task.ReplyTo)
			if err != nil {
				// TODO sleep? break?
				logger.Log(err)
			} else {
				_, err = task.ReplayToAPI.Call(ctx, method, body)
				if err != nil {
					// TODO sleep?
					logger.Log(err)
				}
			}
		}
	}
}
