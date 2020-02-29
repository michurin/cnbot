package api

import (
	"context"
	"encoding/json"
	"github.com/michurin/cnbot/pkg/apirequest"
	"net/http"

	"github.com/michurin/cnbot/pkg/interfaces"
	"github.com/pkg/errors"
)

type API struct {
	apiURLPrefix string
	client       interfaces.HTTPClient
}

func New(client interfaces.HTTPClient, token string) *API {
	return &API{
		apiURLPrefix: "https://api.telegram.org/bot" + token + "/",
		client:       client,
	}
}

func (a *API) Call(ctx context.Context, method string, request apirequest.Request) (json.RawMessage, error) {
	status, body, err := a.client.Do(
		ctx,
		request.Method,
		request.MIME,
		a.apiURLPrefix+method,
		request.Body)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errors.Errorf("HTTP status: %d", status)
	}
	result, err := extractResult(body)
	if err != nil {
		return nil, err
	}
	return result, nil
}
