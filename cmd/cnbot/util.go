package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/michurin/cnbot/pkg/interfaces"
	"github.com/pkg/errors"
)

func sigtermListener(
	rootCtx context.Context,
	logger interfaces.Logger,
) (
	context.Context,
	context.CancelFunc,
) {
	ctx, cancel := context.WithCancel(rootCtx)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-stop
		logger.Log(fmt.Sprintf("Interrupted by signal: %s", s))
		cancel()
	}()
	return ctx, cancel
}

func formatJSON(b []byte) (string, error) {
	x := interface{}(nil)
	err := json.Unmarshal(b, &x)
	if err != nil {
		return "", errors.WithStack(err)
	}
	d, err := json.MarshalIndent(x, "", "    ")
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(d), nil
}

func parseFlags() (string, bool, error) {
	check := flag.Bool("just-check", false, "Just check token. Call getMe method")
	flag.Parse()
	return "bots.ini", *check, nil // TODO filename hardcoded
}
