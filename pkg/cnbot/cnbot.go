package cnbot

import (
	"encoding/json"
	"log"
	"github.com/michurin/cnbot/pkg/cnbot/httpcall"
	"github.com/michurin/cnbot/pkg/cnbot/replay"
	"net/http"
)

const API_CALL_GET_UPDATES = "https://api.telegram.org/bot226015286:AAHzO9VmmKi-_uwZkmbA0DVn03DOpdYYsg4/getUpdates"
const POLLING_TIMOUT = 10

type updatesStruct struct {
	Ok bool
	Result []struct {
		Message struct {
			Chat struct {
				Id int64
			}
			Text string
		}
		UpdateID int64 `json:"update_id"`
	}
}

func PollingLoop(proxy_server string) {
	client_factory := httpcall.ClientFactory{proxy_server}
	apiClient, err := client_factory.NewTAPIClient()
	if err != nil {
		log.Panic(err) // Panic!
	}
	var lastUpdateID int64
	for {
		request_data := map[string]interface{}{"timeout": POLLING_TIMOUT}
		if lastUpdateID > 0 {
			request_data["offset"] = lastUpdateID + 1
		}
		request_body, err := json.Marshal(request_data)
		if err != nil {
			log.Fatal(err) // TODO не так фатально
		}
		log.Print("Longpolling...\n")
		t, err := httpcall.HTTPCall(apiClient, http.MethodPost, API_CALL_GET_UPDATES, request_body)
		if err != nil {
			log.Fatal(err) // TODO не так фатально, можно поспать и продолжать
		}
		resp := updatesStruct{}
		err = json.Unmarshal(t, &resp)
		log.Print(resp)
		for _, message := range resp.Result {
			if message.UpdateID > lastUpdateID {
				lastUpdateID = message.UpdateID
			}
			go replay.Reply(message.Message.Chat.Id, client_factory)
			log.Printf(`update_id=%d text="%s"\n`, message.UpdateID, message.Message.Text)
		}
	}
}