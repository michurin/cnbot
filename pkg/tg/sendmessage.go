package tg

import (
	"encoding/json"
	"errors"

	hps "github.com/michurin/cnbot/pkg/helpers"
)

var markdownMode = hps.StrPtr("MarkdownV2")

type inlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type inlineKeyboardMarkup struct {
	InlineKeyboard [][]inlineKeyboardButton `json:"inline_keyboard"`
}

type sendMessageRequest struct {
	ChatID                int64                 `json:"chat_id"`
	Text                  string                `json:"text"`
	ParseMode             *string               `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool                  `json:"disable_web_page_preview"`
	ReplyMarkup           *inlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func EncodeSendMessage(to int64, text string, isMarkdown bool, markup [][][2]string) (*Request, error) {
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
		ReplyMarkup:           makeInlineKeyboardMarkup(markup),
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

func makeInlineKeyboardMarkup(markup [][][2]string) (mup *inlineKeyboardMarkup) {
	if len(markup) > 0 {
		kbd := [][]inlineKeyboardButton(nil)
		for _, a := range markup {
			x := []inlineKeyboardButton(nil)
			for _, b := range a {
				x = append(x, inlineKeyboardButton{
					Text:         b[1],
					CallbackData: b[0],
				})
			}
			kbd = append(kbd, x)
		}
		mup = &inlineKeyboardMarkup{InlineKeyboard: kbd}
	}
	return
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
