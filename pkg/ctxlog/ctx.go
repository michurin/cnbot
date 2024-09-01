package ctxlog

import (
	"context"
)

type ctxKeyT int

const ctxKey = ctxKeyT(0)

func Add(ctx context.Context, x ...any) context.Context {
	if ox, ok := ctx.Value(ctxKey).([][]any); ok {
		return context.WithValue(ctx, ctxKey, append(ox, x))
	}
	return context.WithValue(ctx, ctxKey, [][]any{x})
}

type PatchAttrs struct {
	attrs [][]any
}

func Patch(ctx context.Context) PatchAttrs {
	if ox, ok := ctx.Value(ctxKey).([][]any); ok {
		return PatchAttrs{attrs: ox}
	}
	return PatchAttrs{attrs: nil}
}

func ApplyPatch(ctx context.Context, p PatchAttrs) context.Context {
	if len(p.attrs) == 0 {
		return ctx
	}
	if ox, ok := ctx.Value(ctxKey).([][]any); ok {
		return context.WithValue(ctx, ctxKey, append(ox, p.attrs...))
	}
	return ctx
}
