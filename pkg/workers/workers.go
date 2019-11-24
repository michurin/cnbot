package workers

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/michurin/cnbot/pkg/api"
	"github.com/michurin/cnbot/pkg/datatype"
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
			if err != nil {
				// TODO sleep? break?
				logger.Log(err)
			}
			fmt.Printf("OUT: %q\n", string(out))
			fmt.Printf("ERR: %+v\n", err)
			fmt.Printf("TASK: %s\n", task.Text) // TODO use logger

			method, body, err := buildRequest(out, task)
			if err != nil {
				// TODO sleep? break?
				logger.Log(err)
			}
			_, err = a.Call(ctx, method, body)
			if err != nil {
				// TODO sleep?
				logger.Log(err)
			}
		}
	}
}

func buildRequest(data []byte, task Task) (string, api.Request, error) {
	var err error
	var body api.Request
	var method string
	imgType := datatype.ImageType(data)
	if imgType != "" {
		method = api.MethodSendPhoto
		body, err = api.EncodeMultipart(task.ReplyTo, data, imgType)
		if err != nil {
			return "", api.Request{}, errors.WithStack(err)
		}
	} else {
		method = api.MethodSendMessage
		body, err = api.EncodeJSON(map[string]interface{}{
			"chat_id": task.ReplyTo,
			"text":    string(data),
		})
		if err != nil {
			return "", api.Request{}, errors.WithStack(err)
		}
	}
	return method, body, nil
}
