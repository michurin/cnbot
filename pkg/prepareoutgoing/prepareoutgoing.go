package prepareoutgoing

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"strconv"

	"github.com/michurin/cnbot/pkg/log"
	"github.com/michurin/cnbot/pkg/sender"
)

var emptyOutgoingData = sender.OutgoingData{}

func PrepareOutgoing(
	log *log.Logger,
	outData []byte,
	chatId int64,
	tips map[string]string,
) (sender.OutgoingData, error) {
	isEmpty, leftIt, rawJSON, isImage, imageType, err := classifyData(outData)
	if err != nil {
		log.Errorf("Classification error: %s", err.Error())
		payload := map[string]interface{}{
			"chat_id": chatId,
			"text":    err.Error(),
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return emptyOutgoingData, err
		}
		return sender.OutgoingData{
			MessageType: "sendMessage",
			Type:        "application/json",
			Body:        body,
		}, nil
	}
	if leftIt {
		return emptyOutgoingData, nil // empty MessageType means left message
	}
	if isEmpty {
		outData = []byte("(no data)")
	}
	if rawJSON {
		return sender.OutgoingData{
			MessageType: "sendMessage",
			Type:        "application/json",
			Body:        outData,
		}, nil
	}
	if isImage {
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		err := w.WriteField("chat_id", strconv.FormatInt(chatId, 10))
		if err != nil {
			return emptyOutgoingData, err
		}
		for _, k := range [...]string{"parse_mode", "disable_notification", "caption"} {
			v, ok := tips[k]
			if ok {
				err := w.WriteField(k, v)
				if err != nil {
					return emptyOutgoingData, err
				}
			}
		}
		fw, err := w.CreateFormFile("photo", "image."+imageType) // TODO: [1] imageType used as extension :-(, [2] Content-Type: application/octet-stream :-(
		if err != nil {
			return emptyOutgoingData, err
		}
		fw.Write(outData)
		w.Close()
		return sender.OutgoingData{
			MessageType: "sendPhoto",
			Type:        w.FormDataContentType(),
			Body:        b.Bytes(),
		}, nil
	} else {
		payload := map[string]interface{}{
			"chat_id": chatId,
			"text":    string(outData),
		}
		for _, k := range [...]string{"parse_mode", "disable_notification"} {
			v, ok := tips[k]
			if ok {
				payload[k] = v
			}
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return emptyOutgoingData, err
		}
		return sender.OutgoingData{
			MessageType: "sendMessage",
			Type:        "application/json",
			Body:        body,
		}, nil
	}
}

func CallbackAnswerOutgoing(callbackQueryId string) (sender.OutgoingData, error) {
	payload := map[string]string{
		"callback_query_id": callbackQueryId,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return emptyOutgoingData, err
	}
	return sender.OutgoingData{
		MessageType: "answerCallbackQuery",
		Type:        "application/json",
		Body:        body,
	}, nil
}
