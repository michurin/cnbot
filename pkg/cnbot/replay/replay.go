package replay

import (
	"log"
	"encoding/json"
	"github.com/michurin/cnbot/pkg/cnbot/httpcall"
	"github.com/michurin/cnbot/pkg/cnbot/getjoke"
	"net/http"
)

const API_CALL_SEND_MESSAGE = "https://api.telegram.org/bot226015286:AAHzO9VmmKi-_uwZkmbA0DVn03DOpdYYsg4/sendMessage"

func Reply(chatId int64, http_client_factory httpcall.ClientFactory) {
	log.Printf("Prepare reply for chat_id = %d", chatId)
	apiClient, err := http_client_factory.NewTAPIClient()
	if err != nil {
		log.Panic(err) // TODO
	}
	joke, err := getjoke.GetJoke()
	if err != nil {
		log.Fatal(err)
	}
	requestBody, err := json.Marshal(map[string]interface{}{"chat_id": chatId, "text": joke})
	if err != nil {
		log.Fatal(err)
	}
	t, err := httpcall.HTTPCall(apiClient, http.MethodPost, API_CALL_SEND_MESSAGE, requestBody)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Result of sendMessage: %s", t)
}
