package tg

import (
	"encoding/json"
	"errors"

	hps "github.com/michurin/cnbot/pkg/helpers"
)

var markdownMode = hps.StrPtr("MarkdownV2")

type sendMessageRequest struct {
	ChatID                int     `json:"chat_id"`
	Text                  string  `json:"text"`
	ParseMode             *string `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool    `json:"disable_web_page_preview"`
}

func EncodeSendMessage(to int, text string, isMarkdown bool) (*Request, error) {
	// In fact text must not be empty and must not be longer then 4K
	// however, 4K *after* parsing. So, we perform this check earlier.
	mode := (*string)(nil)
	if isMarkdown {
		mode = markdownMode
	}
	body, err := json.Marshal(sendMessageRequest{
		ChatID:                to,
		Text:                  text,
		ParseMode:             mode,
		DisableWebPagePreview: true,
	})
	if err != nil {
		return nil, err // in fact, it is the reason for panic
	}
	return &Request{
		Method:      "sendMessage",
		ContentType: "application/json",
		Body:        body,
	}, nil
}

type sendMessageResponse struct {
	Ok bool `json:"ok"`
}

func DecodeSendMessage(body []byte) error {
	data := sendMessageResponse{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	if !data.Ok {
		return errors.New("body is not OK: " + string(body))
	}
	return nil
}
