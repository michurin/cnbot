package calltgapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/michurin/cnbot/pkg/log"
)

func url(token string, method string) string {
	return "https://api.telegram.org/bot" + token + "/" + method
}

func PostBytes(
	log *log.Logger,
	timeout_sec int,
	token string,
	method string,
	data []byte,
	mime string,
	resp interface{},
) error {
	//log.Debugf("Raw send: %v", data)
	timeout := time.Duration(timeout_sec) * time.Second
	response, err := (&http.Client{
		Timeout: timeout,
	}).Post(
		url(token, method),
		mime,
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	//log.Debugf("Raw recv: %v", body)
	err = json.Unmarshal(body, resp)
	if err != nil {
		log.Errorf(
			"Can not parse response [code=%d] [body=%s]",
			response.StatusCode,
			body,
		)
	}
	return err
}
