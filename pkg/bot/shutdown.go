package bot

import (
	"context"
	"os"
	"os/signal"

	hps "github.com/michurin/cnbot/pkg/helpers"
)

// Do not forget defer cancel()
func ShutdownCtx(ctx context.Context, sig ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, sig...)
	go func() {
		sig := <-sigs
		hps.Log(ctx, "Killed by signal", sig.String())
		cancel()
	}()
	return ctx, cancel
}
