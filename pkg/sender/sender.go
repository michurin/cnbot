package sender

import (
	"context"
	"time"

	"github.com/michurin/cnbot/pkg/calltgapi"
	"github.com/michurin/cnbot/pkg/log"
)

type OutgoingData struct {
	MessageType string
	Type        string
	Body        []byte
	Response    chan []byte
}

func Sender(log *log.Logger, token string, toSendQueue <-chan OutgoingData) {
	timeout := time.Duration(10) * time.Second // TODO: send operation timeout, have to be configurable
	for message := range toSendQueue {
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		body, err := calltgapi.PostBytes(
			ctx,
			log,
			token,
			message.MessageType,
			message.Body,
			message.Type,
		)
		if err != nil {
			log.Error(err)
		}
		if message.Response != nil {
			message.Response <- body
			close(message.Response)
		}
		log.Infof("Response: %s", string(body)) // TODO: process it
	}
}
