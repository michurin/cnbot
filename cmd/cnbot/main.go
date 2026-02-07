package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/michurin/cnbot/pkg/app"
	"github.com/michurin/cnbot/pkg/xlog"
)

func main() {
	vFlag := flag.Bool("v", false, "show version and exit")
	dFlag := flag.Bool("d", false, "turn on debugging")
	flag.Parse()
	if vFlag != nil && *vFlag {
		app.ShowVersionInfo()
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	app.SetupLogging(dFlag != nil && *dFlag)
	cfg, tgAPIOrigin, err := app.LoadConfigs(flag.Args()...)
	if err != nil {
		xlog.L(ctx, err)
		return
	}
	err = app.Application(ctx, cfg, tgAPIOrigin, app.MainVersion())
	if err != nil {
		xlog.L(ctx, err)
		return
	}
}
