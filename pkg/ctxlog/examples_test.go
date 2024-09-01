package ctxlog_test

import (
	"context"
	"log/slog"
	"os"

	"github.com/michurin/cnbot/pkg/ctxlog"
)

const thisFileName = "ctxlog/examples_test.go"

var optsNoTimeNoSourceNoLevel = slog.HandlerOptions{
	AddSource: false,
	Level:     nil,
	ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr { // remove time; just to be reproducible
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		return a
	},
}

func ExampleHandler_usualUsecase() {
	// Somewhere you create handler.

	baseHandler := slog.Handler(slog.NewTextHandler(os.Stdout, &optsNoTimeNoSourceNoLevel))

	// You can setup custom attrs for handler. Our wrapper won't manage that attrs.

	baseHandler = baseHandler.WithAttrs([]slog.Attr{slog.Any("app", "one")})

	// Now you are able to setup global logger. You can setup lib-wide or application-wide logger using slog.SetDefault()

	log := slog.New(ctxlog.Handler(baseHandler, thisFileName))

	// You may have a chain of calls in you apps, let's say next two funcs.

	funcContextFreeLogic := func() error {
		return ctxlog.Errorf("initial error")
	}

	funcClient := func(ctx context.Context, arg int) error {
		ctx = ctxlog.Add(ctx, "client", "clientLabel", "arg", arg)
		err := funcContextFreeLogic()
		if err != nil {
			return ctxlog.Errorfx(ctx, "client error: %w", err)
		}
		return nil
	}

	funcHandler := func(ctx context.Context, input int) error {
		ctx = ctxlog.Add(ctx, "component", "handlerLabel")
		err := funcClient(ctx, input)
		if err != nil {
			return ctxlog.Errorfx(ctx, "handler failure: %w", err)
		}
		return nil
	}

	// You instrumentation is able to setup context and call the chain

	ctx := context.Background()

	ctx = ctxlog.Add(ctx, "request_id", "deadbeef")

	err := funcHandler(ctx, -1) // -1 will cause error
	if err != nil {
		log.Error("Error", "key_does_not_matter", err) // key doesn't matter as long as error is wrapped; handler does magic with err
	}

	// output:
	// level=ERROR msg=Error app=one source=ctxlog/examples_test.go:69 err_source=ctxlog/examples_test.go:40 err_msg="handler failure: client error: initial error" request_id=deadbeef component=handlerLabel client=clientLabel arg=-1
}

func ExampleHandler_howGroupsAndAttrsDoing() {
	baseHandler := slog.Handler(slog.NewTextHandler(os.Stdout, &optsNoTimeNoSourceNoLevel))

	log := slog.New(ctxlog.Handler(baseHandler, thisFileName))
	log.Info("Message")
	log.Info("Message-inline-attrs", "P", "Q")
	log.InfoContext(ctxlog.Add(context.Background(), "V", "W"), "Message-1-ctx-attrs")
	log = log.With("X", "Y")
	log.Info("Message-with-attrs")
	log = log.WithGroup("G")
	log.Info("Message-with-group")

	// output:
	// level=INFO msg=Message source=ctxlog/examples_test.go:80
	// level=INFO msg=Message-inline-attrs source=ctxlog/examples_test.go:81 P=Q
	// level=INFO msg=Message-1-ctx-attrs source=ctxlog/examples_test.go:82 V=W
	// level=INFO msg=Message-with-attrs X=Y source=ctxlog/examples_test.go:84
	// level=INFO msg=Message-with-group X=Y G.source=ctxlog/examples_test.go:86
}

func ExampleHandler_indirectContextEnrichment() {
	baseHandler := slog.Handler(slog.NewTextHandler(os.Stdout, &optsNoTimeNoSourceNoLevel))

	log := slog.New(ctxlog.Handler(baseHandler, thisFileName))

	// we populate context somewhere, for example in main.go; we are setting up high level context of logging
	ctx := ctxlog.Add(context.Background(), "component", "A")

	// now we want to create handler or factory, that will be used somewhere else and has to use attrs from this base context
	handler := func(patch ctxlog.PatchAttrs) func(ctx context.Context) { // return handler func
		return func(ctx context.Context) { // ctx here is a handle time context
			ctx = ctxlog.ApplyPatch(ctx, patch) // populate ctx
			log.InfoContext(ctx, "OK")          // here both handle time attr and factory initialization time attrs show up
		}
	}(ctxlog.Patch(ctx))

	// somewhere our handler is called with some low level context
	handlerContext := ctxlog.Add(context.Background(), "request_id", 99)
	handler(handlerContext)

	// output:
	// level=INFO msg=OK source=ctxlog/examples_test.go:108 request_id=99 component=A
}
