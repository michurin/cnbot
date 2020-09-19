package tg

import (
	"bytes"
	"fmt"
	"mime/multipart"
)

func EncodeSendPhoto(to int, tp string, data []byte) (*Request, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	err := w.WriteField("chat_id", fmt.Sprintf("%d", to))
	if err != nil {
		return nil, err
	}
	fw, err := w.CreateFormFile("photo", "image."+tp)
	if err != nil {
		return nil, err
	}
	_, err = fw.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return &Request{
		Method:      "sendPhoto",
		ContentType: w.FormDataContentType(),
		Body:        body.Bytes(),
	}, nil
}

// use DecodeSendMessage to decode response
