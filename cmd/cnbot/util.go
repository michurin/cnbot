package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
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
	token := flag.String("token", "", "You can specify token in command line. However you do not want")
	tokenFile := flag.String("token-file", "", "File name with token")
	flag.Parse()
	if *token == "" {
		if *tokenFile == "" {
			return "", false, errors.New("You have to specify either token or token file.")
		}
		content, err := ioutil.ReadFile(*tokenFile)
		if err != nil {
			return "", false, errors.WithStack(err)
		}
		return strings.TrimSpace(string(content)), *check, nil
	}
	return *token, *check, nil
}
