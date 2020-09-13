package tg

import (
	"bytes"
	"fmt"
	"mime/multipart"
)

func EncodeSendPhoto(to int, tp string, data []byte) Request {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	err := w.WriteField("chat_id", fmt.Sprintf("%d", to))
	if err != nil {
		return Request{Error: err}
	}
	fw, err := w.CreateFormFile("photo", "image."+tp)
	if err != nil {
		return Request{Error: err}
	}
	_, err = fw.Write(data)
	if err != nil {
		return Request{Error: err}
	}
	err = w.Close()
	if err != nil {
		return Request{Error: err}
	}
	return Request{
		Method:      "sendPhoto",
		ContentType: w.FormDataContentType(),
		Body:        body.Bytes(),
	}
}
