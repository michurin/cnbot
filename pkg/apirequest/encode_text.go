package apirequest

import "github.com/pkg/errors"

func TextMessage(text string, to int, isMarkdown bool) (Request, error) {
	if text == "" {
		return Request{}, errors.New("empty text message. I skip it")
	}
	data := map[string]interface{}{
		"chat_id": to,
		"text":    text,
	}
	if isMarkdown {
		data["parse_mode"] = "markdown"
	}
	req, err := EncodeJSON(data)
	if err != nil {
		return Request{}, err
	}
	return req, nil
}
