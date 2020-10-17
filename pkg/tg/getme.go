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
		FirstName               string `json:"first_name"`
		Username                string `json:"username"`
		ID                      int    `json:"id"`
		IsBot                   bool   `json:"is_bot"`
		CanJoinGroups           bool   `json:"can_join_groups"`
		CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
		SupportsInlineQueries   bool   `json:"supports_inline_queries"`
	}
}

func DecodeGetMe(body []byte) (
	id int,
	username string,
	firstname string,
	canJoinGrp bool,
	canReadAllGrpMsg bool,
	supportInline bool,
	err error,
) {
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
	r := data.Result
	id = r.ID
	username = r.Username
	firstname = r.FirstName
	canJoinGrp = r.CanJoinGroups
	canReadAllGrpMsg = r.CanReadAllGroupMessages
	supportInline = r.SupportsInlineQueries
	return
}
