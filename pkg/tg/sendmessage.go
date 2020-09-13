package tg

import "encoding/json"

type sendMessageRequest struct {
	ChatID                int    `json:"chat_id"`
	Text                  string `json:"text"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

func EncodeSendMessage(to int, text string) Request {
	body, err := json.Marshal(sendMessageRequest{
		ChatID:                to,
		Text:                  text,
		DisableWebPagePreview: true,
	})
	if err != nil {
		return Request{Error: err}
	}
	return Request{
		Method: "sendMessage",
		Body:   body,
	}
}
