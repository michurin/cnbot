package tg

import "github.com/michurin/cnbot/pkg/helpers"

const tgPrefixUrl = "https://api.telegram.org/bot"

type Request struct {
	Method      string
	ContentType string
	Body        []byte
}

func Encode(token string, r *Request) helpers.Request {
	url := tgPrefixUrl + token + "/" + r.Method
	if r.Body == nil {
		return helpers.Request{
			Method: "GET",
			URL:    url,
		}
	}
	return helpers.Request{
		Method:      "POST",
		URL:         url,
		ContentType: r.ContentType,
		Body:        r.Body,
	}
}
