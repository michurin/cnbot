package api

import (
	"bytes"
	"mime/multipart"
	"strconv"

	"github.com/pkg/errors"
)

func EncodeMultipart(chatID int, data []byte) (Request, error) {

	typeExtension := "png" // TODO

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	err := w.WriteField("chat_id", strconv.Itoa(chatID))
	if err != nil {
		return Request{}, errors.WithStack(err)
	}

	fw, err := w.CreateFormFile("photo", "image."+typeExtension)
	if err != nil {
		return Request{}, errors.WithStack(err)
	}
	n, err := fw.Write(data)
	if err != nil {
		return Request{}, errors.WithStack(err)
	}
	if n != len(data) {
		return Request{}, errors.New("Not all data has been written")
	}

	w.Close()

	return Request{
		Method: "POST",
		MIME:   w.FormDataContentType(),
		Body:   b.Bytes(),
	}, nil
}
