package xbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/michurin/cnbot/pkg/ctxlog"
	"github.com/michurin/cnbot/pkg/xlog"
)

type Request struct {
	Method      string
	ContentType string
	Body        []byte
}

func RequestStruct(method string, x any) (*Request, error) {
	d, err := json.Marshal(x)
	if err != nil {
		return nil, err
	}
	return &Request{
		Method:      method,
		ContentType: "application/json",
		Body:        d,
	}, nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func RequestFromBinary(data []byte, userID int64) (*Request, error) {
	contentType := http.DetectContentType(data)
	switch {
	case strings.HasPrefix(contentType, "text/"):
		if !utf8.Valid(data) {
			return nil, ctxlog.Errorf("invalid utf8")
		}
		if bytes.HasPrefix(data, []byte("%!PRE\n")) {
			croped := bytes.TrimSpace(data[6:]) // 6 is len of prefix
			if len(croped) == 0 {
				return RequestStruct("sendMessage", map[string]any{
					"chat_id":  userID,
					"text":     "-", // Forcing empty preformated text to "-"
					"entities": []any{map[string]any{"type": "pre", "offset": 0, "length": 1}},
				})
			}
			str := string(croped)
			r := []rune(str)
			if len(r) > 4096 { // In this place length is in terms of Unicode chars
				return reqMultipart("sendDocument", userID, "document", croped, "formatted_text", "text/plain") // TODO nice to add caption... with caption_entities?
			}
			return RequestStruct("sendMessage", map[string]any{
				"chat_id": userID,
				"text":    str,
				"entities": []any{
					map[string]any{
						"type":   "pre",
						"offset": 0,
						"length": len(utf16.Encode(r)), // According Telegram's API, in this place the length is in terms of UTF-16
					},
				},
			})
		}
		data = bytes.TrimSpace(data)
		if len(data) == 0 {
			return nil, ctxlog.Errorf("zero length data: skipping response")
		}
		str := string(data)
		r := []rune(str)
		if len(r) > 4096 { // In this place length is in terms of Unicode chars
			return reqMultipart("sendDocument", userID, "document", data, "message", "text/plain") // TODO nice to add caption
		}
		return RequestStruct("sendMessage", map[string]any{"chat_id": userID, "text": str})
	case strings.HasPrefix(contentType, "image/"): // TODO to limit image formats
		return reqMultipart("sendPhoto", userID, "photo", data, "image", contentType)
	case strings.HasPrefix(contentType, "video/"): // TODO to limit video formats
		return reqMultipart("sendVideo", userID, "video", data, "video", contentType)
	case strings.HasPrefix(contentType, "audio/"): // it seems application/ogg is not fully supported; it requires OPUS encoding
		return reqMultipart("sendAudio", userID, "audio", data, "audio", contentType)
	default: // TODO hmm... application/* and font/*
		xlog.L(context.TODO(), fmt.Sprintf("Fallback to multipart from %q", contentType)) // TODO no context here
		return reqMultipart("sendDocument", userID, "document", data, "document", contentType)
	}
}

func fext(ctype string) string {
	// mime.ExtensionsByType can return several extension sorted alphabetical.
	// We are trying to find most common extension, by comparing with last
	// part of mime type of data.
	prefExt, _, _ := mime.ParseMediaType(ctype)
	idx := strings.LastIndex(prefExt, "/")
	if idx >= 0 {
		prefExt = "." + prefExt[idx+1:]
	}
	exts, _ := mime.ExtensionsByType(ctype)
	if len(exts) == 0 {
		return ".dat"
	}
	for _, e := range exts { // find preferable extension, if any
		if e == prefExt {
			return e
		}
	}
	e := exts[0]
	if e == ".asc" {
		return ".txt" // looks more reasonable for text/plain
	}
	return e
}

func reqMultipart(method string, to int64, fieldname string, data []byte, filename string, ctype string) (*Request, error) { // TODO legacy
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	err := w.WriteField("chat_id", strconv.FormatInt(to, 10))
	if err != nil {
		return nil, err
	}
	fw, err := w.CreatePart(textproto.MIMEHeader{ // see implementation of CreateFormFile
		"Content-Disposition": {fmt.Sprintf(
			`form-data; name="%s"; filename="%s"`,
			quoteEscaper.Replace(fieldname),
			quoteEscaper.Replace(filename+fext(ctype)),
		)},
		"Content-Type": {ctype},
	})
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
		Method:      method,
		ContentType: w.FormDataContentType(),
		Body:        body.Bytes(),
	}, nil
}
