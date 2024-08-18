package ctxlog

import (
	"context"
	"errors"
	"fmt"
	"runtime"
)

type wrapError struct {
	err   error
	attrs [][]any
	pc    uintptr
}

func Errorf(format string, a ...any) error {
	err := fmt.Errorf(format, a...) //nolint:goerr113
	x := new(wrapError)
	if errors.As(err, &x) { // do not wrap twice
		return &wrapError{err: err, attrs: x.attrs, pc: x.pc} // new message (Error()), but attrs and pc from wrapped error
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	return &wrapError{err: err, attrs: nil, pc: pcs[0]}
}

func Errorfx(ctx context.Context, format string, a ...any) error {
	err := fmt.Errorf(format, a...) //nolint:goerr113
	x := new(wrapError)
	if errors.As(err, &x) { // already wrapped
		attrs := x.attrs  // however it is possible we do not have any attrs yet
		if attrs == nil { // in this case we are fill attrs
			if a, ok := ctx.Value(ctxKey).([][]any); ok {
				attrs = a
			}
		}
		return &wrapError{err: err, attrs: attrs, pc: x.pc} // new message (Error()), but attrs and pc from wrapped error
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	attrs := [][]any(nil)
	if a, ok := ctx.Value(ctxKey).([][]any); ok {
		attrs = a
	}
	return &wrapError{err: err, attrs: attrs, pc: pcs[0]}
}

func (e *wrapError) Error() string {
	return e.err.Error()
}

func (e *wrapError) Unwrap() error {
	return e.err
}
