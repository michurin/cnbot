package tg

import (
	hps "github.com/michurin/cnbot/pkg/helpers"
)

const tgPrefixURL = "https://api.telegram.org/bot"

type Request struct {
	Method      string
	ContentType string
	Body        []byte
}

func Encode(token string, r *Request) hps.Request {
	url := tgPrefixURL + token + "/" + r.Method
	if r.Body == nil {
		return hps.Request{
			Method: "GET",
			URL:    url,
		}
	}
	return hps.Request{
		Method:      "POST",
		URL:         url,
		ContentType: r.ContentType,
		Body:        r.Body,
	}
}
