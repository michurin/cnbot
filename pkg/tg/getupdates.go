package tg

import (
	"encoding/json"
	"errors"
)

type getUpdateRequest struct {
	Offset         *int64   `json:"offset,omitempty"`
	Timeout        int      `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

func EncodeGetUpdates(offset int64, timeout int) (*Request, error) {
	if timeout <= 0 {
		return nil, errors.New("timeout have to be greater then zero")
	}
	r := getUpdateRequest{
		Timeout:        timeout,
		AllowedUpdates: []string{"message", "callback_query"},
	}
	if offset != 0 {
		r.Offset = &offset
	}
	body, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return &Request{
		Method:      "getUpdates",
		ContentType: "application/json",
		Body:        body,
	}, nil
}

type Message struct {
	BotName         string
	Text            string
	FromID          int64
	FromFirstName   string
	ChatID          int64
	CallbackID      string
	UpdateMessageID int64
	SideType        string
	SideID          int64
	SideName        string
}

type user struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	IsBot     bool   `json:"is_bot"`
}

type chat struct {
	ID    int64  `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"` // optional
}

type contact struct {
	PhoneNumber string  `json:"phone_number"` // TODO not used yet
	FirstName   string  `json:"first_name"`
	LastName    *string `json:"last_name"` // TODO not used yet
	UserID      *int64  `json:"user_id"`
	Vcard       *string `json:"vcard"` // TODO not used yet
}

type message struct {
	MessageID       int64    `json:"message_id"`
	Text            string   `json:"text"`
	From            user     `json:"from"`
	Chat            chat     `json:"chat"`
	ForwardFrom     *user    `json:"forward_from"`
	ForwardFromChat *chat    `json:"forward_from_chat"`
	Contact         *contact `json:"contact"`
}

type callbackQuery struct {
	ID      string  `json:"id"`
	From    user    `json:"from"`
	Message message `json:"message"` // TODO it is optional
	Data    string  `json:"data"`
}

type getUpdateResponse struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID      int64          `json:"update_id"`
		Message       *message       `json:"message"`
		CallbackQuery *callbackQuery `json:"callback_query"`
	} `json:"result"`
}

// Slightly magically cares about offset. It return previous offset if no
// messages or errors.
func DecodeGetUpdates(body []byte, offset int64, botName string) ([]Message, int64, error) {
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
		if e.Message != nil {
			msg := e.Message
			m[i].FromID = msg.From.ID
			m[i].FromFirstName = msg.From.FirstName
			m[i].ChatID = msg.Chat.ID
			m[i].Text = msg.Text
			sideType, sideUserID, sideUserName := extractSideUser(msg)
			m[i].SideID = sideUserID
			m[i].SideName = sideUserName
			m[i].SideType = sideType
		} else if e.CallbackQuery != nil {
			cb := e.CallbackQuery
			m[i].FromID = cb.From.ID
			m[i].FromFirstName = cb.From.FirstName
			m[i].ChatID = cb.Message.Chat.ID
			m[i].CallbackID = cb.ID
			m[i].UpdateMessageID = cb.Message.MessageID
			m[i].Text = cb.Data
		}
		if e.UpdateID > u {
			u = e.UpdateID
		}
	}
	return m, u + 1, nil
}

func extractSideUser(msg *message) (tp string, id int64, name string) {
	if msg.ForwardFromChat != nil {
		tp = msg.ForwardFromChat.Type
		id = msg.ForwardFromChat.ID
		name = msg.ForwardFromChat.Title
		return
	}
	if msg.ForwardFrom != nil {
		if msg.ForwardFrom.IsBot {
			tp = "bot"
		} else {
			tp = "user"
		}
		id = msg.ForwardFrom.ID
		name = msg.ForwardFrom.FirstName
		return
	}
	if msg.Contact != nil && msg.Contact.UserID != nil {
		tp = "contact"
		id = *msg.Contact.UserID
		name = msg.Contact.FirstName
		return
	}
	return
}
