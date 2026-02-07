package xlog

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"sync/atomic"
	"time"
)

var defaultLogger atomic.Pointer[slog.Logger]

func init() { //nolint:gochecknoinits
	defaultLogger.Store(slog.Default())
}

// SetDefault mimics slog.SetDefault
func SetDefault(l *slog.Logger) {
	defaultLogger.Store(l)
}

// L is botwide logging function, it could be private or internal
// it is best with ctxlog.Handler
func L(ctx context.Context, a any) {
	log(ctx, a, slog.LevelInfo)
}

// D is a debug-variant of L
func D(ctx context.Context, a any) {
	log(ctx, a, slog.LevelDebug)
}

func log(ctx context.Context, a any, normalLogLevel slog.Level) {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // skip
	r := slog.Record{}
	switch v := a.(type) {
	case error:
		r = slog.NewRecord(time.Now(), slog.LevelError, v.Error(), pcs[0])
		r.Add("raw_error", v) // if v is wrapped error, the key will be skipped in ctxlog.Handler, value will be interpreted and split to several key-value pairs
	case string:
		r = slog.NewRecord(time.Now(), normalLogLevel, v, pcs[0])
	case []byte:
		r = slog.NewRecord(time.Now(), normalLogLevel, safeString(v), pcs[0])
	case nil:
		r = slog.NewRecord(time.Now(), normalLogLevel, "<nil>", pcs[0])
	default:
		r = slog.NewRecord(time.Now(), slog.LevelWarn, fmt.Sprintf("%[1]T: %#[1]v", a), pcs[0])
	}
	h := defaultLogger.Load().Handler()
	if h.Enabled(ctx, r.Level) { // it seems it is reason to use high level slog.Logger
		_ = h.Handle(ctx, r)
	}
}
