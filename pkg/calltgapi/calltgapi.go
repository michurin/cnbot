package calltgapi

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/michurin/cnbot/pkg/log"
)

func url(token string, method string) string {
	return "https://api.telegram.org/bot" + token + "/" + method
}

func PostBytes(
	ctx context.Context,
	log *log.Logger,
	token string,
	method string,
	data []byte,
	mime string,
) ([]byte, error) {
	//log.Debugf("Raw send: %v", data)
	req, err := http.NewRequest(http.MethodPost, url(token, method), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mime)
	req.Header.Set("User-Agent", "CNBot (https://github.com/michurin/cnbot/)") // TODO to config, version?
	response, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	//log.Debugf("Raw recv: %v", body)
	return body, nil
}
