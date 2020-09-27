package helpers

import (
	"context"
	"time"
)

func Sleep(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(duration):
	}
}
