package tg

import (
	"encoding/json"
)

type editMessageRequest struct {
	ChatID                int64                 `json:"chat_id"`
	MessageID             int64                 `json:"message_id"`
	Text                  string                `json:"text"`
	ParseMode             *string               `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool                  `json:"disable_web_page_preview"`
	ReplyMarkup           *inlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func EncodeEditMessage(chatID, messageID int64, text string, isMarkdown bool, markup [][][2]string) (*Request, error) {
	// In fact text must not be empty and must not be longer then 4K
	// however, 4K *after* parsing. So, we perform this check earlier.
	mode := (*string)(nil)
	if isMarkdown {
		mode = markdownMode
	}
	body, err := json.Marshal(editMessageRequest{
		ChatID:                chatID,
		MessageID:             messageID,
		Text:                  text,
		ParseMode:             mode,
		DisableWebPagePreview: true,
		ReplyMarkup:           makeInlineKeyboardMarkup(markup),
	})
	if err != nil {
		return nil, err // in fact, it is the reason for panic
	}
	return &Request{
		Method:      "editMessageText",
		ContentType: "application/json",
		Body:        body,
	}, nil
}

// use DecodeSendMessage to decode response
