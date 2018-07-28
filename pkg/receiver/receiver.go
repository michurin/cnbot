package receiver

import (
	"encoding/json"

	"github.com/michurin/cnbot/pkg/calltgapi"
	"github.com/michurin/cnbot/pkg/log"
)

type TUpdateChat struct {
	FirstName string `json:"first_name"`
	Id        int64  `json:"id"`
	LastName  string `json:"last_name"`
	Type      string `json:"type"`
	Username  string `json:"username"`
}

type TUpdateFrom struct {
	FirstName    string `json:"first_name"`
	Id           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	LanguageCode string `json:"language_code"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
}

type TUpdateMessage struct { // TODO: only this struct have to be public, isn't it?
	MessageId int64       `json:"message_id"`
	Date      int64       `json:"date"`
	Text      string      `json:"text"`
	Chat      TUpdateChat `json:"chat"`
	From      TUpdateFrom `json:"from"`
}

type TUpdateResult struct {
	UpdateId int64          `json:"update_id"`
	Message  TUpdateMessage `json:"message"`
}

type TUpdates struct {
	Ok     bool            `json:"ok"`
	Result []TUpdateResult `json:"result"`
}

type TUpdatesRequest struct {
	Offset         *int64   `json:"offset,omitempty"`
	Limit          *int64   `json:"limit,omitempty"`
	Timeout        *int     `json:"timeout,omitempty"`
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
}

func RunPollingLoop(
	log *log.Logger,
	apiTimeout int,
	token string,
	messages chan<- TUpdateMessage,
) {
	var offset int64
	var offsetPtr *int64 // have to be nil on first iteration
	var updates *TUpdates
	for {
		body, err := calltgapi.PostBytes(
			log,
			apiTimeout+10,
			token,
			"getUpdates",
			mustMarshal(
				log,
				TUpdatesRequest{Offset: offsetPtr, Timeout: &apiTimeout},
			),
			"application/json",
		)
		if err != nil {
			log.Error(err)
			errorSleep(log, 10)
			continue
		}
		updates = &TUpdates{}
		err = json.Unmarshal(body, updates)
		if err != nil {
			log.Error(err)
			errorSleep(log, 10)
			continue
		}
		log.Debugf("Update %#v", updates)
		if !updates.Ok {
			log.Error("ERROR: Update is not Ok")
			errorSleep(log, 20)
			continue
		}
		for _, part := range updates.Result {
			if offsetPtr == nil {
				offsetPtr = &offset
			}
			if offset <= part.UpdateId {
				offset = part.UpdateId + 1
			}
			log.Info(part.Message)
			messages <- part.Message
		}
		log.Debugf("offset = %d", offset)
	}
}
