package interfaces

import "context"

type HTTPClient interface {
	Do(ctx context.Context, method string, contentType string, url string, body []byte) (respStatus int, repBody []byte, err error)
}
