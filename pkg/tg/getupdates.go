package tg

import (
	"encoding/json"
	"errors"
)

type getUpdateRequest struct {
	Offset         *int     `json:"offset,omitempty"`
	Timeout        int      `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

func EncodeGetUpdates(offset int, timeout int) Request {
	if timeout <= 0 {
		return Request{Error: errors.New("timeout have to be greater then zero")}
	}
	r := getUpdateRequest{
		Timeout:        timeout,
		AllowedUpdates: []string{"message"},
	}
	if offset != 0 {
		r.Offset = &offset
	}
	body, err := json.Marshal(r)
	if err != nil {
		return Request{Error: err}
	}
	return Request{
		Method:      "getUpdates",
		ContentType: "application/json",
		Body:        body,
	}
}

type Message struct {
	BotName string
	Text    string
	FromID  int
}

type getUpdateResponse struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			Text string `json:"text"`
			From struct {
				ID int `json:"id"`
			} `json:"from"`
		} `json:"message"`
	} `json:"result"`
}

// Slightly magically cares about offset. It return previous offset if no
// messages or errors.
func DecodeGetUpdate(body []byte, offset int, botName string) ([]Message, int, error) {
	data := getUpdateResponse{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, offset, err
	}
	if !data.Ok {
		return nil, offset, errors.New("body is not OK: " + string(body))
	}
	if len(data.Result) == 0 {
		return nil, offset, nil
	}
	m := make([]Message, len(data.Result))
	u := data.Result[0].UpdateID
	for i, e := range data.Result {
		m[i].BotName = botName
		m[i].Text = e.Message.Text
		m[i].FromID = e.Message.From.ID
		if e.UpdateID > u {
			u = e.UpdateID
		}
	}
	return m, u + 1, nil
}
