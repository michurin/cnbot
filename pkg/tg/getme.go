package tg

import (
	"encoding/json"
	"errors"
)

func EncodeGetMe() (request *Request) {
	return &Request{Method: "getMe"}
}

type getMeResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		ID        int    `json:"id"`
		IsBot     bool   `json:"is_bot"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
	}
}

func DecodeGetMe(body []byte) (id int, username string, err error) {
	data := getMeResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}
	if !data.Ok {
		err = errors.New("body is not OK: " + string(body))
		return
	}
	if !data.Result.IsBot {
		err = errors.New("bot is not bot: " + string(body))
		return
	}
	id = data.Result.ID
	username = data.Result.Username
	return
}
