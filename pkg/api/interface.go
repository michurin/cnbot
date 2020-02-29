package api

import (
	"context"
	"encoding/json"
)

type Interface interface {
	Call(ctx context.Context, method string, request Request) (json.RawMessage, error)
}

