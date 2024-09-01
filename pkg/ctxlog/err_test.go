package ctxlog_test

import (
	"errors"
	"testing"

	"github.com/michurin/cnbot/pkg/ctxlog"
)

func TestErrWrap(t *testing.T) {
	specificErr := errors.New("x")
	err := ctxlog.Errorf("err: %w", specificErr)
	t.Log(errors.Is(err, specificErr))
}
