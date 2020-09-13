package tg

import "github.com/michurin/cnbot/pkg/helpers"

const tgPrefixUrl = "https://api.telegram.org/bot"

type Request struct {
	Method      string
	ContentType string
	Body        []byte
	Error       error // we include error to simplify usage of encoders
}

func Encode(token string, r Request) (helpers.Request, error) {
	if r.Error != nil {
		return helpers.Request{}, r.Error
	}
	url := tgPrefixUrl + token + "/" + r.Method
	if r.Body == nil {
		return helpers.Request{
			Method: "GET",
			URL:    url,
		}, nil
	}
	return helpers.Request{
		Method:      "POST",
		URL:         url,
		ContentType: r.ContentType,
		Body:        r.Body,
	}, nil
}
