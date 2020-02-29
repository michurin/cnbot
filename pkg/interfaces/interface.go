package interfaces

import (
	"context"
	"encoding/json"
	"github.com/michurin/cnbot/pkg/apirequest"
)

type Interface interface {
	Call(ctx context.Context, method string, request apirequest.Request) (json.RawMessage, error)
}

