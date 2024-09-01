package xbot

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/michurin/cnbot/pkg/ctxlog"
	"github.com/michurin/cnbot/pkg/xlog"
)

type Client interface {
	Do(r *http.Request) (*http.Response, error)
}

type Bot struct {
	APIOrigin string // injection to be testable
	Token     string
	Client    Client // injection to be observable // TODO move all logging into client middleware?
}

func (b *Bot) API(ctx context.Context, request *Request) ([]byte, error) {
	ctx = xlog.API(ctx, request.Method, request.ContentType)
	err := error(nil)
	req := (*http.Request)(nil)
	resp := (*http.Response)(nil)
	respCode := 0
	data := []byte(nil)
	defer func() {
		msg := any(nil)
		if err != nil {
			msg = err
		} else {
			msg = "ok"
		}
		xlog.L(xlog.Status(xlog.Request(xlog.Response(ctx, data), request.Body), respCode), msg)
	}()
	reqURL := b.APIOrigin + "/bot" + b.Token + "/" + request.Method
	req, err = http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(request.Body))
	if err != nil {
		return nil, ctxlog.Errorfx(ctx, "request constructor: %w", err)
	}
	req.Header.Set("Content-Type", request.ContentType)
	resp, err = b.Client.Do(req)
	if err != nil {
		return nil, ctxlog.Errorfx(ctx, "client: %w", err)
	}
	respCode = resp.StatusCode
	defer resp.Body.Close()           // we are skipping error here
	data, err = io.ReadAll(resp.Body) // we have to read and close Body even for non-200 responses
	if err != nil {
		return nil, ctxlog.Errorfx(ctx, "reading: %w", err)
	}
	return data, nil
}

func (b *Bot) Download(ctx context.Context, path string, stream io.Writer) error {
	ctx = xlog.API(ctx, "x-download", "x")
	err := error(nil)
	defer func() {
		xlog.L(xlog.Path(ctx, path), err)
	}()
	reqURL := b.APIOrigin + "/file/bot" + b.Token + "/" + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return ctxlog.Errorfx(ctx, "request constructor: %w", err)
	}
	resp, err := b.Client.Do(req)
	if err != nil {
		return ctxlog.Errorfx(ctx, "client: %w", err)
	}
	defer resp.Body.Close() // we are skipping error here
	_, err = io.Copy(stream, resp.Body)
	if err != nil {
		return ctxlog.Errorfx(ctx, "coping: %w", err)
	}
	return nil
}
