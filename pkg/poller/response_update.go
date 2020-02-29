package poller

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Update struct {
	UpdateID int `json:"update_id"`
	Message  *struct {
		Text string `json:"text"`
		Date int    `json:"date"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			UserName  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		From struct {
			ID           int    `json:"id"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			UserName     string `json:"username"`
			Type         string `json:"type"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"from"`
	} `json:"message"`
}

func ParseUpdates(r json.RawMessage) ([]Update, error) {
	x := []Update(nil)
	err := json.Unmarshal(r, &x)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return x, nil
}
