package receiver

import (
	"encoding/json"

	"github.com/michurin/cnbot/pkg/log"
)

func mustMarshal(log *log.Logger, data interface{}) []byte {
	body, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return body
}
