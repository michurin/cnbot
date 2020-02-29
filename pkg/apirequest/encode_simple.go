package apirequest

import (
	"strings"

	"github.com/michurin/cnbot/pkg/datatype"
	"github.com/pkg/errors"
)

// It is not true encoder. It's multiencoder.
// I'm not sure here is a best place for it.
// TODO move it closer to server.go?
func SimpleRequest(data []byte, replyTo int) (string, Request, error) {
	var err error
	var req Request
	var method string
	imgType := datatype.ImageType(data)
	if imgType != "" {
		method = MethodSendPhoto
		req, err = EncodeMultipart(replyTo, data, imgType, "", false)
		if err != nil {
			return "", Request{}, errors.WithStack(err)
		}
	} else {
		method = MethodSendMessage
		text := strings.TrimSpace(string(data))
		if strings.HasPrefix(text, "```") {
			req, err = TextMessage(string(data), replyTo, true)
			if err != nil {
				return "", Request{}, errors.WithStack(err)
			}
		} else {
			req, err = TextMessage(string(data), replyTo, false)
			if err != nil {
				return "", Request{}, errors.WithStack(err)
			}
		}
	}
	return method, req, nil
}
