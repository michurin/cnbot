package xlog

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"unicode/utf8"

	"github.com/michurin/cnbot/pkg/ctxlog"
)

func Bot(ctx context.Context, name string) context.Context {
	return ctxlog.Add(ctx, slog.String("bot", name))
}

func Comp(ctx context.Context, name string) context.Context {
	return ctxlog.Add(ctx, slog.String("comp", name))
}

func API(ctx context.Context, method, contentType string) context.Context {
	return ctxlog.Add(ctx, slog.String("api", method), slog.String("content_type", contentType))
}

func Status(ctx context.Context, status int) context.Context {
	return ctxlog.Add(ctx, slog.Int("status", status))
}

func Request(ctx context.Context, request []byte) context.Context {
	return ctxlog.Add(ctx, slog.String("request", safeString(request)))
}

func Response(ctx context.Context, response []byte) context.Context {
	return ctxlog.Add(ctx, slog.String("response", safeString(response)))
}

func Path(ctx context.Context, path string) context.Context {
	return ctxlog.Add(ctx, slog.String("path", path))
}

func User(ctx context.Context, user int64) context.Context {
	return ctxlog.Add(ctx, slog.Int64("user", user))
}

func Pid(ctx context.Context, pid int) context.Context {
	return ctxlog.Add(ctx, slog.Int("pid", pid))
}

func safeString(x []byte) string { // TODO move to field-specific place (does not exist yet)
	idx := bytes.IndexByte(x, '\n')
	if idx >= 0 {
		x = x[:idx] // side effect prone
	}
	if utf8.Valid(x) {
		return string(x)
	}
	return fmt.Sprintf("%q", x)
}
