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
	for message := range toSendQueue {
		body, err := calltgapi.PostBytes(
			log,
			10, // TODO: send operation timeout, have to be configurable
			token,
			message.MessageType,
			message.Body,
			message.Type,
		)
		if err != nil {
			log.Error(err)
		}
		log.Infof("Response: %s", string(body)) // TODO: process it
	}
}
