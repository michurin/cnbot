package sender

import (
	"github.com/michurin/cnbot/pkg/calltgapi"
	"github.com/michurin/cnbot/pkg/log"
)

type OutgoingData struct {
	MessageType string
	Type        string
	Body        []byte
}

func Sender(log *log.Logger, token string, toSendQueue <-chan OutgoingData) {
	for {
		message := <-toSendQueue
		resp := map[string]interface{}{}
		err := calltgapi.PostBytes(
			log,
			10, // TODO: send operation timeout, have to be configurable
			token,
			message.MessageType,
			message.Body,
			message.Type,
			&resp,
		)
		if err != nil {
			log.Error(err)
		}
		resp_ok, ok := resp["ok"].(bool)
		if !ok {
			log.Errorf("Resp error: %v", resp)
		} else if !resp_ok {
			log.Errorf("Resp not ok: %v", resp)
		}
	}
}
