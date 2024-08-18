package app

import (
	"context"
	"net/http"
	"path"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/michurin/cnbot/pkg/ctxlog"
	"github.com/michurin/cnbot/pkg/xbot"
	"github.com/michurin/cnbot/pkg/xcfg"
	"github.com/michurin/cnbot/pkg/xctrl"
	"github.com/michurin/cnbot/pkg/xlog"
	"github.com/michurin/cnbot/pkg/xloop"
	"github.com/michurin/cnbot/pkg/xproc"
)

func bot(ctx context.Context, eg *errgroup.Group, cfg xcfg.Config, tgAPIOrigin, build string) {
	bot := &xbot.Bot{
		APIOrigin: tgAPIOrigin,
		Token:     cfg.Token,
		Client:    http.DefaultClient,
	}

	envCommon := []string{"tg_x_ctrl_addr=" + cfg.ControlAddr, "tg_x_build=" + build}

	command := &xproc.Cmd{
		InterruptDelay: 10 * time.Second,
		KillDelay:      10 * time.Second,
		Env:            envCommon,
		Command:        cfg.Script,
		ConfigFileDir:  cfg.ConfigFileDir,
	}

	commandLong := &xproc.Cmd{
		InterruptDelay: 10 * time.Minute,
		KillDelay:      10 * time.Minute,
		Env:            envCommon,
		Command:        cfg.LongRunningScript,
		ConfigFileDir:  path.Dir(cfg.LongRunningScript),
	}

	eg.Go(func() error {
		err := xloop.Loop(xlog.Comp(ctx, "loop"), bot, command)
		if err != nil {
			return ctxlog.Errorfx(ctx, "polling loop: %w", err)
		}
		return nil
	})

	server := &http.Server{
		Addr: cfg.ControlAddr,
		Handler: xctrl.Handler(
			bot,
			commandLong,
			ctxlog.Patch(xlog.Comp(ctx, "ctrl")),
		),
		ReadHeaderTimeout: time.Minute, // slowloris attack protection
	}
	eg.Go(func() error {
		<-ctx.Done()
		cx, stop := context.WithTimeout(context.Background(), time.Second)
		defer stop()
		return server.Shutdown(cx) //nolint:contextcheck
	})

	eg.Go(func() error {
		err := server.ListenAndServe()
		if err != nil {
			return ctxlog.Errorfx(ctx, "control server listener: %w", err)
		}
		return nil
	})
}

func Application(rootCtx context.Context, bots map[string]xcfg.Config, tgAPIOrigin, build string) error {
	if len(bots) == 0 {
		return ctxlog.Errorfx(rootCtx, "there is no configuration")
	}
	eg, ctx := errgroup.WithContext(rootCtx)
	for name, cfg := range bots {
		bot(xlog.Bot(ctx, name), eg, cfg, tgAPIOrigin, build)
	}
	xlog.L(ctx, "Run. Build="+build)
	return eg.Wait()
}
