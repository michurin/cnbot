package app

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	"github.com/michurin/cnbot/pkg/ctxlog"
	"github.com/michurin/cnbot/pkg/xlog"
)

// logHandler implements interface slog.Handler
// it is drop-in replacement for slog.NewTextHandler, but more human friendly
type logHandler struct{}

func (logHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (logHandler) Handle(_ context.Context, r slog.Record) error {
	kv := map[string]any{} // not thread safe, however r.Attrs works consequently
	r.Attrs(func(a slog.Attr) bool {
		kv[a.Key] = a.Value.Any()
		return true
	})
	std := ""                                                    // std attributes
	for _, a := range []string{"bot", "comp", "api", "source"} { // order significant
		if v, ok := kv[a]; ok {
			std = std + " [" + v.(string) + "]" //nolint:forcetypeassert // we use typed helples to enrich context with all this values
			delete(kv, a)
		}
	}
	ekeys := []string(nil) // extra keys
	for k := range kv {
		ekeys = append(ekeys, k)
	}
	sort.Strings(ekeys)
	nstd := ""
	for _, a := range ekeys {
		nstd += fmt.Sprintf(" %s=%v", a, kv[a])
	}
	fmt.Printf("%s [%s]%s%s %s\n", r.Time.Format("2006-01-02 15:04:05"), r.Level.String(), std, nstd, r.Message)
	return nil
}

func (h logHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	panic("NOT IMPLEMENTED")
}

func (logHandler) WithGroup(_ string) slog.Handler {
	panic("NOT IMPLEMENTED")
}

func SetupLogging() {
	l := slog.New(ctxlog.Handler(logHandler{}, "app/log.go"))
	xlog.SetDefault(l)
}
