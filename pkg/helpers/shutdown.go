package helpers

import (
	"context"
	"os"
	"os/signal"

	"github.com/michurin/minlog"
)

// Do not forget defer cancel()
func ShutdownCtx(ctx context.Context, sig ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, sig...)
	go func() {
		sig := <-sigs
		minlog.Log(ctx, "Killed by signal", sig.String())
		cancel()
	}()
	return ctx, cancel
}
