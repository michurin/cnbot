package api

import (
	"context"
	"encoding/json"
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

func (a *API) JSON(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	var err error
	var body []byte
	if params != nil {
		body, err = json.Marshal(params)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	status, body, err := a.client.Do(ctx, "POST", "application/json", a.apiURLPrefix+method, body)
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
