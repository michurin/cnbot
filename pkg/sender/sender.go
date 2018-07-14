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
		err := calltgapi.PostBytes(log, token, message.MessageType, message.Body, message.Type, &resp)
		if err != nil {
			log.Error(err)
		}
		if !resp["ok"].(bool) {
			log.Errorf("Resp error: %v", resp)
		}
	}
}
