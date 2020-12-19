package tg

import (
	"encoding/json"
)

type answerCallbackQueryRequest struct {
	CallbackQueryID string `json:"callback_query_id"`
}

func EncodeAnswerCallbackQuery(callbackID string) (*Request, error) {
	body, err := json.Marshal(answerCallbackQueryRequest{
		CallbackQueryID: callbackID,
	})
	if err != nil {
		return nil, err // in fact, it is the reason for panic
	}
	return &Request{
		Method:      "answerCallbackQuery",
		ContentType: "application/json",
		Body:        body,
	}, nil
}
