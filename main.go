package main

import (
	"context"
	"os"
	"syscall"

	"github.com/michurin/cnbot/pkg/bot"
)

func main() {
	ctx, cancel := bot.ShutdownCtx(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()
	bot.Run(ctx)
}
