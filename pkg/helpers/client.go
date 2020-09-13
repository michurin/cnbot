package helpers

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
)

type Request struct {
	Method      string
	URL         string
	ContentType string
	Body        []byte
}

func Do(ctx context.Context, req Request) ([]byte, error) {
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, bytes.NewReader(req.Body))
	if err != nil {
		Log(ctx, err)
		return nil, err
	}
	if req.ContentType != "" {
		httpReq.Header.Set("Content-Type", req.ContentType)
	}
	Log(ctx, httpReq.Method, httpReq.URL.String(), req.Body)
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		Log(ctx, err)
		return nil, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			Log(ctx, err)
		}
	}()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log(ctx, err)
		return nil, err
	}
	Log(ctx, resp.StatusCode, respBody)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("not 200")
	}
	return respBody, nil
}
