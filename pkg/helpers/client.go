package helpers

import (
	"bytes"
	"context"
	"errors"
	"io"
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
	if resp != nil && resp.Body != nil {
		defer func() {
			// in go 1.5 we have to read all the rest data from buffer if any to be able to reuse connection
			_, err := io.Copy(ioutil.Discard, resp.Body)
			if err != nil {
				Log(ctx, err)
			}
			err = resp.Body.Close()
			if err != nil {
				Log(ctx, err)
			}
		}()
	}
	if err != nil {
		Log(ctx, err)
		return nil, err
	}
	if resp.Body == nil {
		return nil, nil
	}
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
