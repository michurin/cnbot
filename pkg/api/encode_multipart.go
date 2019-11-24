package api

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"strconv"

	"github.com/pkg/errors"
)

func EncodeMultipart(chatID int, data []byte, typeExtension string) (Request, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	err := w.WriteField("chat_id", strconv.Itoa(chatID))
	if err != nil {
		return Request{}, errors.WithStack(err)
	}

	contentDesc := fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "photo", "image."+typeExtension)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", contentDesc)
	h.Set("Content-Type", "image/"+typeExtension)
	fw, err := w.CreatePart(h)
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

	err = w.Close()
	if err != nil {
		return Request{}, errors.WithStack(err)
	}

	return Request{
		Method: "POST",
		MIME:   w.FormDataContentType(),
		Body:   b.Bytes(),
	}, nil
}
