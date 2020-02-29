package server

import (
	"fmt"
	"github.com/michurin/cnbot/pkg/api"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/michurin/cnbot/pkg/datatype"

	"github.com/michurin/cnbot/pkg/apirequest"
	"github.com/michurin/cnbot/pkg/interfaces"

	"github.com/pkg/errors"
)

type HTTPHandler struct {
	logger interfaces.Logger
	api    *api.API
}

func New(logger interfaces.Logger, api *api.API) *HTTPHandler {
	return &HTTPHandler{
		logger: logger,
		api:    api,
	}
}

// echo -n PART | curl -F 'do=@-' 'http://localhost:9999/ok'
// echo -n PART | curl -d '@-' 'http://localhost:9999/ok'
// TODO:
// - use logger
// - call API
// - write response to client
// - check http method
func (h HTTPHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Printf("SERVER %s [NOLOG!]", req.URL)              // TODO
	log.Printf("Method: %s", req.Method)                   // TODO
	log.Printf("MIME: %s", req.Header.Get("Content-Type")) // TODO
	err := h.serve(req)
	if err != nil {
		h.logger.Log(err)
		_, _ = resp.Write([]byte(err.Error())) // TODO
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusOK)
	_, _ = resp.Write([]byte("")) // TODO
}

func (h HTTPHandler) serve(req *http.Request) error {
	mediaType, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		return errors.WithStack(err)
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		return h.multi(req, params["boundary"])
	}
	return h.simple(req)
}

func (h HTTPHandler) simple(req *http.Request) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errors.WithStack(err)
	}
	replyToStr := req.URL.Query().Get("to")
	if replyToStr == "" {
		return errors.New("empty parameter 'to'")
	}
	replyTo, err := strconv.Atoi(replyToStr)
	if err != nil {
		return errors.WithStack(err)
	}

	method, apiReq, err := apirequest.SimpleRequest(body, replyTo)
	if err != nil {
		return err
	}
	resp, err := h.api.Call(req.Context(), method, apiReq)
	if err != nil {
		return err
	}
	_ = resp // TODO return it
	return nil
}

func (h HTTPHandler) multi(req *http.Request, boundary string) error {
	replyToStr := req.URL.Query().Get("to")
	if replyToStr == "" {
		return errors.New("empty parameter 'to'")
	}
	replyTo, err := strconv.Atoi(replyToStr)
	if err != nil {
		return errors.WithStack(err)
	}
	mr := multipart.NewReader(req.Body, boundary)
	parts := map[string][]byte{}
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break // When there are no more parts, the error io.EOF is returned.
		}
		if err != nil {
			return errors.WithStack(err)
		}
		slurp, err := ioutil.ReadAll(p)
		if err != nil {
			return errors.WithStack(err)
		}
		mediaType, params, err := mime.ParseMediaType(p.Header.Get("Content-Disposition"))
		if err != nil {
			return errors.WithStack(err)
		}
		fmt.Printf("Part [%s] %q: %q\n", mediaType, params["name"], slurp)
		parts[params["name"]] = slurp
	}

	text := string(parts["text"])
	_, isMarkdown := parts["markdown"]
	if data, ok := parts["photo"]; ok {
		imgType := datatype.ImageType(data)
		if imgType == "" {
			return errors.Errorf("can not detect image type: %q", data[:10]) // TODO RANGE!!!
		}
		apiReq, err := apirequest.EncodeMultipart(replyTo, data, imgType, text, isMarkdown)
		if err != nil {
			return err
		}
		resp, err := h.api.Call(req.Context(), apirequest.MethodSendPhoto, apiReq)
		_ = resp // TODO
	} else {
		apiReq, err := apirequest.TextMessage(text, replyTo, isMarkdown)
		if err != nil {
			return err
		}
		resp, err := h.api.Call(req.Context(), apirequest.MethodSendMessage, apiReq)
		_ = resp // TODO
	}
	return nil
}
