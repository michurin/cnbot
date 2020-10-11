package tg

import (
	"encoding/json"
	"errors"
)

func EncodeGetWebhookInfo() (request *Request) {
	return &Request{Method: "getWebhookInfo"}
}

type getWebhookInfoResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		URL string `json:"url"`
	}
}

func DecodeGetWebhookInfo(body []byte) (url string, err error) {
	data := getWebhookInfoResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}
	if !data.Ok {
		err = errors.New("body is not OK: " + string(body))
		return
	}
	url = data.Result.URL
	return
}
