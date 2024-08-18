package ctxlog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
)

type handler struct {
	pfx  string
	next slog.Handler
}

func Handler(next slog.Handler, sfx string) slog.Handler {
	pfx := ""
	_, file, _, ok := runtime.Caller(1)
	if ok {
		p, ok := strings.CutSuffix(file, sfx)
		if ok {
			pfx = p
		}
	}
	return &handler{pfx: pfx, next: next}
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *handler) fname(pc uintptr) string {
	f, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	file, _ := strings.CutPrefix(f.File, h.pfx)
	return fmt.Sprintf("%s:%d", file, f.Line)
}

func (h *handler) Handle(ctx context.Context, ro slog.Record) error {
	// Slightly overcomplicated approach. So we want to:
	// - save attrs from handler as is
	// - find and skip wrapped errors
	// - deduplicate keys, that come from ctx and err
	// We want to manage !BADKEY naturally, preserve order of keys and
	// certain priority.
	// Beware: we does not see here attrs of next handler, so we do not consider and process them.
	r := slog.NewRecord(ro.Time, ro.Level, ro.Message, ro.PC) // result record; we are enable duplicates here
	if ro.PC != 0 {
		r.AddAttrs(slog.Any(slog.SourceKey, h.fname(ro.PC)))
	}
	lastErr := (*wrapError)(nil) // not thread safe code; h.Attrs doing consequently
	ro.Attrs(func(a slog.Attr) bool {
		v := a.Value.Any()
		if err, ok := v.(error); ok {
			e := new(wrapError)
			if errors.As(err, &e) {
				lastErr = e
				return true
			}
		}
		r.AddAttrs(a)
		return true
	})
	ri := slog.NewRecord(ro.Time, ro.Level, ro.Message, ro.PC) // interim record; to manage !BADKEY etc in native way
	if aa, ok := ctx.Value(ctxKey).([][]any); ok {
		for _, a := range aa {
			ri.Add(a...)
		}
	}
	if lastErr != nil {
		r.AddAttrs(slog.Any("err_source", h.fname(lastErr.pc)), slog.Any("err_msg", lastErr.Error()))
		for _, a := range lastErr.attrs {
			ri.Add(a...)
		}
	}
	idx := map[string]int{}
	i := 0
	ri.Attrs(func(a slog.Attr) bool {
		idx[a.Key] = i
		i++
		return true
	})
	i = 0
	ri.Attrs(func(a slog.Attr) bool {
		if idx[a.Key] == i {
			r.AddAttrs(a)
		}
		i++
		return true
	})
	return h.next.Handle(ctx, r) //nolint:wrapcheck // this error will be ignored at log/slog/logger.log()
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{
		pfx:  h.pfx,
		next: h.next.WithAttrs(attrs),
	}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		pfx:  h.pfx,
		next: h.next.WithGroup(name),
	}
}
