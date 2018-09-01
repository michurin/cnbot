package httpif

import (
	"encoding/json"
	"errors"

	"github.com/michurin/cnbot/pkg/receiver"
)

type TResponse struct {
	Ok     bool            `json:"ok"`
	Result json.RawMessage `json:"result"` // On deleteMessage result=true and can NOT be parsed :-(
}

var errorMessageNotOk = errors.New("Message is not OK")

func replyToMessageId(data []byte) (int64, error) {
	hlMess := new(TResponse)
	err := json.Unmarshal(data, hlMess)
	if err != nil {
		return 0, err
	}
	if !hlMess.Ok {
		return 0, errorMessageNotOk
	}
	mess := new(receiver.TUpdateMessage)
	err = json.Unmarshal(hlMess.Result, mess)
	if err != nil {
		return 0, err
	}
	return mess.MessageId, nil
}
