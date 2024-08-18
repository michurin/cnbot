package xlog_test

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/michurin/cnbot/pkg/ctxlog"
	"github.com/michurin/cnbot/pkg/xlog"
)

// WARNING: L is global and this tests can ruin other tests
// --------------------------------------------------------

var optToBeReproducible = &slog.HandlerOptions{
	ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		return a
	},
}

func ExampleL() {
	xlog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, optToBeReproducible)))
	ctx := context.Background()
	xlog.L(ctx, "ok")
	xlog.L(ctx, errors.New("err"))
	// Output:
	// level=INFO msg=ok
	// level=ERROR msg=err raw_error=err
}

func ExampleL_withCtxLogHandler() {
	xlog.SetDefault(slog.New(ctxlog.Handler(slog.NewTextHandler(os.Stdout, optToBeReproducible), "log_test.go")))
	ctx := context.Background()
	ctx = ctxlog.Add(ctx, "label", "A")
	xlog.L(ctx, "ok")
	err := func(ctx context.Context) error { // some random function
		ctx = ctxlog.Add(ctx, "scope", "S")
		err := errors.New("err")
		return ctxlog.Errorfx(ctx, "error: %w", err)
	}(ctx)
	xlog.L(ctx, err) // we do not have scope=S in this ctx, however, we can see it in logs
	// Output:
	// level=INFO msg=ok source=log_test.go:39 label=A
	// level=ERROR msg="error: err" source=log_test.go:45 err_source=log_test.go:43 err_msg="error: err" label=A scope=S
}
