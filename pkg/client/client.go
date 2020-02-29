// Weary simple http client. Even simpler than standard http.Client
// I use it everywhere in this library as interface. So you are free
// to make your own implementation or simple wrapper as I do in logging wrapper.

package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type HTTPClient struct {
	client http.Client
}

func New(client http.Client) *HTTPClient {
	return &HTTPClient{client: client}
}

func (h *HTTPClient) Do(ctx context.Context, method string, contentType string, url string, body []byte) (int, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return 0, nil, errors.WithStack(err)
	}
	req.Header.Set("Content-Type", contentType)
	resp, err := h.client.Do(req)
	if err != nil {
		return 0, nil, errors.WithStack(err)
	}
	defer func() { _ = resp.Body.Close() }()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, errors.WithStack(err)
	}
	return resp.StatusCode, respBody, nil
}
