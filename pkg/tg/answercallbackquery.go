package tg

import (
	"encoding/json"
)

type answerCallbackQueryRequest struct {
	CallbackQueryID string  `json:"callback_query_id"`
	Text            *string `json:"text,omitempty"`
}

func EncodeAnswerCallbackQuery(callbackID, text string) (*Request, error) {
	if callbackID == "" {
		return nil, nil
	}
	a := answerCallbackQueryRequest{
		CallbackQueryID: callbackID,
	}
	if text != "" {
		a.Text = &text
	}
	body, err := json.Marshal(a)
	if err != nil {
		return nil, err // in fact, it is the reason for panic
	}
	return &Request{
		Method:      "answerCallbackQuery",
		ContentType: "application/json",
		Body:        body,
	}, nil
}
