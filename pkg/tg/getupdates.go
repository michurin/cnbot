package tg

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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
		AllowedUpdates: []string{"message", "edited_message", "callback_query"},
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
	FromLastName    string
	FromUsername    string
	FromIsBot       bool
	FromLanguage    string
	// locations are used in environment variables only
	// so we can use strings to keep it
	LocationLongitude            string // floats here, empty string means it is no location in message
	LocationLatitude             string
	LocationHorizontalAccuracy   string // strings instead *float to
	LocationLivePeriod           string // convey nil pointers as
	LocationHeading              string // empty strings
	LocationProximityAlertRadius string
}

type user struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	IsBot        bool   `json:"is_bot"`
	LastName     string `json:"last_name"`     // optional
	Username     string `json:"username"`      // optional
	LanguageCode string `json:"language_code"` // optional
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

type location struct {
	Longitude            float64  `json:"longitude"`
	Latitude             float64  `json:"latitude"`
	HorizontalAccuracy   *float64 `json:"horizontal_accuracy"`
	LivePeriod           *int64   `json:"live_period"`
	Heading              *int64   `json:"heading"`
	ProximityAlertRadius *int64   `json:"proximity_alert_radius"`
}

type message struct {
	MessageID       int64     `json:"message_id"`
	Text            string    `json:"text"` // if fact, text is optional
	From            user      `json:"from"`
	Chat            chat      `json:"chat"`
	ForwardFrom     *user     `json:"forward_from"`
	ForwardFromChat *chat     `json:"forward_from_chat"`
	Contact         *contact  `json:"contact"`
	Location        *location `json:"location"`
}

type callbackQuery struct {
	ID      string   `json:"id"`
	From    user     `json:"from"`
	Message *message `json:"message"`
	Data    string   `json:"data"`
}

type getUpdateResponse struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID      int64          `json:"update_id"`
		Message       *message       `json:"message"`
		EditedMessage *message       `json:"edited_message"`
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
		var from user
		var chatID int64
		switch {
		case e.Message != nil:
			msg := e.Message
			from = msg.From
			chatID = msg.Chat.ID
			m[i].Text = msg.Text
			m[i].SideType, m[i].SideID, m[i].SideName = extractSideUser(msg)
			m[i].LocationLongitude, m[i].LocationLatitude, m[i].LocationHorizontalAccuracy,
				m[i].LocationLivePeriod, m[i].LocationHeading,
				m[i].LocationProximityAlertRadius = extractLocation(msg.Location)
		case e.EditedMessage != nil:
			msg := e.EditedMessage
			from = msg.From
			chatID = msg.Chat.ID
			m[i].Text = msg.Text
			m[i].SideType, m[i].SideID, m[i].SideName = extractSideUser(msg)
			m[i].LocationLongitude, m[i].LocationLatitude, m[i].LocationHorizontalAccuracy,
				m[i].LocationLivePeriod, m[i].LocationHeading,
				m[i].LocationProximityAlertRadius = extractLocation(msg.Location)
		case e.CallbackQuery != nil:
			cb := e.CallbackQuery
			from = cb.From
			m[i].CallbackID = cb.ID
			m[i].Text = cb.Data
			if cb.Message != nil {
				msg := cb.Message
				m[i].UpdateMessageID = msg.MessageID
				chatID = msg.Chat.ID
			} else {
				chatID = from.ID // fallback if no message
			}
		}
		m[i].FromID = from.ID
		m[i].FromFirstName = from.FirstName
		m[i].ChatID = chatID
		m[i].FromLastName = from.LastName
		m[i].FromUsername = from.Username
		m[i].FromIsBot = from.IsBot
		m[i].FromLanguage = from.LanguageCode
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

func extractLocation(loc *location) (lon, lat, acc, period, heading, alertRad string) {
	if loc == nil {
		return
	}
	lon = ftoa(loc.Longitude)
	lat = ftoa(loc.Latitude)
	acc = fPtrToA(loc.HorizontalAccuracy)
	period = iPtrToA(loc.LivePeriod)
	heading = iPtrToA(loc.Heading)
	alertRad = iPtrToA(loc.ProximityAlertRadius)
	return
}

func ftoa(x float64) string {
	return fmt.Sprintf("%f", x)
}

func fPtrToA(x *float64) string {
	if x != nil {
		return fmt.Sprintf("%f", *x)
	}
	return ""
}

func iPtrToA(x *int64) string {
	if x != nil {
		return strconv.FormatInt(*x, 10)
	}
	return ""
}
