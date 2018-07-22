package prepareoutgoing

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"strconv"

	"github.com/michurin/cnbot/pkg/log"
	"github.com/michurin/cnbot/pkg/sender"
)

func classifyData(data []byte) (
	isEmpty bool,
	leftIt bool,
	rawJSON bool,
	isImage bool,
	imageType string,
) {
	// TODO: naive implementation. Details: https://en.wikipedia.org/wiki/List_of_file_signatures
	if len(data) == 0 { // data is always trimed
		isEmpty = true
	} else if bytes.HasPrefix(data, []byte{'{'}) { // TODO: check tail?
		rawJSON = true
	} else if bytes.HasPrefix(data, []byte{'.'}) { // TODO: check tail?
		leftIt = true
	} else if bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}) { // TODO: naive
		isImage = true
		imageType = "jpeg"
	} else if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		isImage = true
		imageType = "png"
	} else if bytes.HasPrefix(data, []byte{0x47, 0x49, 0x46, 0x38}) {
		isImage = true
		imageType = "gif"
	}
	return
}

func PrepareOutgoing(log *log.Logger, outData []byte, chatId int64, tips map[string]string) sender.OutgoingData {
	isEmpty, leftIt, rawJSON, isImage, imageType := classifyData(outData)
	if leftIt {
		log.Info("Left message")
		return sender.OutgoingData{} // empty MessageType means left message
	}
	if isEmpty {
		outData = []byte("(no data)")
	}
	if rawJSON {
		return sender.OutgoingData{
			MessageType: "sendMessage",
			Type:        "application/json",
			Body:        outData,
		}
	}
	if isImage {
		log.Infof("TRY TO SEND IMAGE TYPE %v", imageType)
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		err := w.WriteField("chat_id", strconv.FormatInt(chatId, 10))
		if err != nil {
			panic(err) // TODO
		}
		for _, k := range [...]string{"parse_mode", "disable_notification", "caption"} {
			v, ok := tips[k]
			if ok {
				err := w.WriteField(k, v)
				if err != nil {
					panic(err) // TODO
				}
			}
		}
		fw, err := w.CreateFormFile("photo", "image."+imageType) // TODO: [1] imageType used as extension :-(, [2] Content-Type: application/octet-stream :-(
		if err != nil {
			panic(err) // TODO
		}
		fw.Write(outData)
		w.Close()
		return sender.OutgoingData{
			MessageType: "sendPhoto",
			Type:        w.FormDataContentType(),
			Body:        b.Bytes(),
		}
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
		_ = err // TODO!!!
		return sender.OutgoingData{
			MessageType: "sendMessage",
			Type:        "application/json",
			Body:        body,
		}
	}
}

func CallbackAnswerOutgoing(callbackQueryId string) sender.OutgoingData {
	payload := map[string]string{
		"callback_query_id": callbackQueryId,
	}
	body, err := json.Marshal(payload)
	_ = err // TODO!!!
	return sender.OutgoingData{
		MessageType: "answerCallbackQuery",
		Type:        "application/json",
		Body:        body,
	}
}
