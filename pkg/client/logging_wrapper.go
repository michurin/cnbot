package client

import (
	"context"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/michurin/cnbot/pkg/interfaces"
)

type LogHTTPClient struct {
	next   interfaces.HTTPClient
	logger interfaces.Logger
}

func WithLogging(next interfaces.HTTPClient, logger interfaces.Logger) *LogHTTPClient {
	return &LogHTTPClient{
		next:   next,
		logger: logger,
	}
}

func (h *LogHTTPClient) Do(ctx context.Context, method string, contentType string, url string, body []byte) (int, []byte, error) {
	respBody := []byte(nil)
	respStatus := 0
	err := error(nil)
	start := time.Now()
	defer func() {
		if err == nil {
			h.logger.Log(fmt.Sprintf(
				"[%dms] %s %s %s -> %d %s",
				time.Since(start)/time.Millisecond,
				method, url, strBody(body), respStatus, strBody(respBody)))
		} else {
			h.logger.Log(fmt.Sprintf(
				"[%dms] %s %s %s -> ERROR: %+v",
				time.Since(start)/time.Millisecond,
				method, url, strBody(body), err))
		}
	}()
	respStatus, respBody, err = h.next.Do(ctx, method, contentType, url, body)
	return respStatus, respBody, err
}

func strBody(b []byte) string {
	if b == nil {
		return "<nil_body>"
	}
	if len(b) == 0 {
		return "<empty_body>"
	}
	if utf8.Valid(b) {
		return string(b)
	}
	s := fmt.Sprintf("%q", b)
	if len(s) > 400 {
		s = s[:400] + "..."
	}
	return s
}
