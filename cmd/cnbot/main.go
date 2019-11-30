package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/michurin/cnbot/pkg/api"
	"github.com/michurin/cnbot/pkg/cfg"
	"github.com/michurin/cnbot/pkg/client"
	"github.com/michurin/cnbot/pkg/execute"
	"github.com/michurin/cnbot/pkg/interfaces"
	"github.com/michurin/cnbot/pkg/log"
	"github.com/michurin/cnbot/pkg/poller"
	"github.com/michurin/cnbot/pkg/processors"
	"github.com/michurin/cnbot/pkg/server"
	"github.com/michurin/cnbot/pkg/workers"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger := log.New()

	rootCtx := context.Background()

	breakableCtx, cancel := sigtermListener(rootCtx, logger)
	defer cancel()

	configFileNane, check, err := parseFlags()
	if err != nil {
		logger.Log(err)
		return
	}

	configs, err := cfg.Read(configFileNane, logger)
	if err != nil {
		logger.Log(err)
		return
	}

	pollingClient := client.WithLogging(client.New(http.Client{Timeout: 60 * time.Second}), logger)
	sendClient := client.WithLogging(client.New(http.Client{Timeout: 20 * time.Second}), logger)

	if check {
		err := checkBots(breakableCtx, sendClient, configs)
		if err != nil {
			logger.Log(err)
		}
	} else {
		err := launch(breakableCtx, logger, configs, pollingClient, sendClient)
		if err != nil {
			logger.Log(err)
		}
	}
}

func checkBots(ctx context.Context, client interfaces.HTTPClient, configs []cfg.BotConfig) error {
	for _, conf := range configs {
		fmt.Printf("Bot configuration:\n%s\n", conf)
		a := api.New(client, conf.Token)
		body, err := a.Call(ctx, api.MethodGetMe, api.EncodeEmpty())
		if err != nil {
			return err
		}
		str, err := formatJSON(body)
		if err != nil {
			return err
		}
		fmt.Printf("Bot info:\n%s\n", str)
	}
	return nil
}

func launch(
	ctx context.Context,
	logger interfaces.Logger,
	configs []cfg.BotConfig,
	pollingClient interfaces.HTTPClient,
	sendClient interfaces.HTTPClient,
) error {
	taskQueue := make(chan workers.Task, 100)

	eg, ctx := errgroup.WithContext(ctx)

	multiAPI := map[string]*api.API{} // we don't protect this map by lock, we only read it in go-routines

	mx := http.NewServeMux()

	for _, conf := range configs {
		conf := conf
		eg.Go(func() error {
			return poller.Poller(
				ctx,
				logger,
				conf.Name,
				api.New(pollingClient, conf.Token),
				conf.Script,
				processors.Safe, // TODO make it configurable
				taskQueue)
		})
		a := api.New(sendClient, conf.Token)
		multiAPI[conf.Name] = a
		mx.Handle("/"+conf.Name, server.New(logger, a))
	}

	bindAddress := "127.0.0.1:9999" // TODO make it configurable
	executor := execute.New(logger, []string{
		"PATH=/bin:/usr/bin:/usr/local/bin", // TODO configurable
		"T_BIND=" + bindAddress,
	})
	eg.Go(func() error {
		return workers.QueueProcessor(ctx, logger, executor, taskQueue, multiAPI)
	})
	eg.Go(func() error {
		return serve(ctx, logger, mx, bindAddress)
	})

	return eg.Wait()
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
