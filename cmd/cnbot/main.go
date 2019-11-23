package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/michurin/cnbot/pkg/api"
	"github.com/michurin/cnbot/pkg/client"
	"github.com/michurin/cnbot/pkg/execute"
	"github.com/michurin/cnbot/pkg/interfaces"
	"github.com/michurin/cnbot/pkg/log"
	"github.com/michurin/cnbot/pkg/poller"
	"github.com/michurin/cnbot/pkg/server"
	"github.com/michurin/cnbot/pkg/workers"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func main() {
	token, check, err := parseFlags()

	logger := log.New()

	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		return
	}

	if check {
		c := client.WithLogging(client.New(http.Client{Timeout: 5 * time.Second}), logger)
		a := api.New(c, token)
		body, err := a.JSON(context.Background(), api.MethodGetMe, nil)
		if err != nil {
			fmt.Printf("Error: %+v\n", err)
			return
		}
		str, err := formatJSON(body)
		if err != nil {
			fmt.Printf("Telegram response: %s\n", string(body))
			return
		}
		fmt.Printf("Bot info:\n%s\n", str)
		return
	}

	pollingClient := client.WithLogging(client.New(http.Client{Timeout: 60 * time.Second}), logger)

	rootCtx := context.Background()

	breakableCtx, cancel := sigtermListener(rootCtx, logger)
	defer cancel()

	executor := execute.New(logger)

	taskQueue := make(chan workers.Task, 100)
	eg, ctx := errgroup.WithContext(breakableCtx)
	eg.Go(func() error { return poller.Poller(ctx, logger, api.New(pollingClient, token), taskQueue) })
	eg.Go(func() error {
		return workers.QueueProcessor(ctx, logger, executor, taskQueue, api.New(pollingClient, token))
	})
	eg.Go(func() error { return serve(ctx, logger, server.HTTPHandler{}, "0.0.0.0:9999") })
	err = eg.Wait()

	_ = err // TODO
}

func serve(
	ctx context.Context,
	logger interfaces.Logger,
	handler http.Handler,
	addr string,
) error {
	logger.Log("Server started")
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	serverErrsChannel := make(chan error)
	go func() {
		err := srv.ListenAndServe()
		serverErrsChannel <- errors.WithStack(err)
	}()
	select {
	case err := <-serverErrsChannel:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
