package app

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/michurin/cnbot/pkg/ctxlog"
	"github.com/michurin/cnbot/pkg/xlog"
)

// logHandler implements interface slog.Handler
// it is drop-in replacement for slog.NewTextHandler, but more human friendly
type logHandler struct {
	level slog.Level
}

func (h logHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (logHandler) Handle(_ context.Context, r slog.Record) error {
	kv := map[string]any{} // not thread safe, however r.Attrs works consequently
	r.Attrs(func(a slog.Attr) bool {
		kv[a.Key] = a.Value.Any()
		return true
	})
	std := make([]string, 0, 4)                                  // std attributes
	for _, a := range []string{"bot", "comp", "api", "source"} { // order significant
		if v, ok := kv[a]; ok {
			std = append(std, " ["+v.(string)+"]") //nolint:forcetypeassert // we use typed helples to enrich context with all this values
			delete(kv, a)
		}
	}
	ekeys := make([]string, 0, len(kv)) // extra keys
	for k := range kv {
		ekeys = append(ekeys, k)
	}
	sort.Strings(ekeys)
	nstd := make([]string, len(ekeys))
	for i, a := range ekeys {
		nstd[i] = fmt.Sprintf(" %s=%v", a, kv[a])
	}
	fmt.Printf(
		"%s [%s]%s%s %s\n",
		r.Time.Format("2006-01-02 15:04:05"),
		r.Level.String(),
		strings.Join(std, ""),
		strings.Join(nstd, ""),
		r.Message)
	return nil
}

func (h logHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	panic("NOT IMPLEMENTED")
}

func (logHandler) WithGroup(_ string) slog.Handler {
	panic("NOT IMPLEMENTED")
}

func SetupLogging(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	l := slog.New(ctxlog.Handler(logHandler{level: level}, "app/log.go"))
	xlog.SetDefault(l)
}
