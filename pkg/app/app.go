package app

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/michurin/cnbot/pkg/ctxlog"
	"github.com/michurin/cnbot/pkg/xbot"
	"github.com/michurin/cnbot/pkg/xctrl"
	"github.com/michurin/cnbot/pkg/xlog"
	"github.com/michurin/cnbot/pkg/xloop"
	"github.com/michurin/cnbot/pkg/xproc"
)

func bot(ctx context.Context, eg *errgroup.Group, cfg BotConfig, tgAPIOrigin, build string) error {
	bot := &xbot.Bot{
		APIOrigin: tgAPIOrigin,
		Token:     cfg.Token,
		Client:    http.DefaultClient,
	}

	envCommon := []string{"tg_x_ctrl_addr=" + cfg.ControlAddr, "tg_x_build=" + build}

	cfgDir, err := filepath.Abs(cfg.ConfigFileDir) // if dir is "", it uses CWD
	if err != nil {
		return ctxlog.Errorfx(ctx, "invalid config dir: %w", err)
	}

	// caution:
	// do not return errors and do not interrupt flow after goroutines running
	// it will lead to goroutine leaking

	command := &xproc.Cmd{
		InterruptDelay: 10 * time.Second,
		KillDelay:      10 * time.Second,
		Env:            envCommon,
		Command:        absPathToBinary(cfgDir, cfg.Script),
	}

	commandLong := &xproc.Cmd{
		InterruptDelay: 10 * time.Minute,
		KillDelay:      10 * time.Minute,
		Env:            envCommon,
		Command:        absPathToBinary(cfgDir, cfg.LongRunningScript),
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

	return nil
}

func absPathToBinary(cwd, b string) string {
	// This function assumes that cwd is absolute. It returns absolute path to executable file.
	// case: is it absolute
	if path.IsAbs(b) {
		return b
	}
	// case: is it regular binary, accessible via $path
	if !strings.Contains(b, string(os.PathSeparator)) {
		x, err := exec.LookPath(b)
		if err == nil {
			return x
		}
	}
	return path.Join(cwd, b)
}

type BotConfig struct {
	ControlAddr       string
	Token             string
	Script            string
	LongRunningScript string
	ConfigFileDir     string
}

func Application(rootCtx context.Context, bots map[string]BotConfig, tgAPIOrigin, build string) error {
	if len(bots) == 0 {
		return ctxlog.Errorfx(rootCtx, "there is no configuration")
	}
	eg, ctx := errgroup.WithContext(rootCtx)
	for name, cfg := range bots {
		err := bot(xlog.Bot(ctx, name), eg, cfg, tgAPIOrigin, build)
		if err != nil {
			return err
		}
	}
	xlog.L(ctx, "Run. Build="+build)
	return eg.Wait()
}
