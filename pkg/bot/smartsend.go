package bot

import (
	"context"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

func buildRequest(destUser, callbackMessageID int64, callbackID string, stdout []byte) (req, cbReq *tg.Request, err error) {
	imgExt := hps.ImageType(stdout)
	if imgExt != "" {
		req, err = tg.EncodeSendPhoto(destUser, imgExt, stdout)
		if err != nil {
			return
		}
		cbReq, err = tg.EncodeAnswerCallbackQuery(callbackID, "")
		return
	}
	ignore, msg, isMarkdown, forUpdate, markup, callbackText, err := hps.MessageType(stdout)
	if err != nil {
		return
	}
	cbReq, err = tg.EncodeAnswerCallbackQuery(callbackID, callbackText)
	if err != nil {
		return
	}
	if ignore {
		return
	}
	if forUpdate && callbackMessageID > 0 {
		req, err = tg.EncodeEditMessage(destUser, callbackMessageID, msg, isMarkdown, markup)
		if err != nil {
			return
		}
		return
	}
	req, err = tg.EncodeSendMessage(destUser, msg, isMarkdown, markup)
	if err != nil {
		return
	}
	return
}

func SmartSend(ctx context.Context, token, callbackID string, destUser, callbackMessageID int64, stdout []byte) error {
	req, cbReq, err := buildRequest(destUser, callbackMessageID, callbackID, stdout)
	if err != nil {
		hps.Log(ctx, err)
		return err
	}
	err = sendRequest(ctx, token, cbReq)
	if err != nil {
		hps.Log(ctx, err)
		return err
	}
	err = sendRequest(ctx, token, req)
	if err != nil {
		hps.Log(ctx, err)
		return err
	}
	return nil
}

func sendRequest(ctx context.Context, token string, req *tg.Request) error {
	if req == nil { // message to be ignored
		return nil
	}
	body, err := hps.Do(ctx, tg.Encode(token, req))
	if err != nil {
		return err
	}
	err = tg.DecodeSendMessage(body)
	if err != nil {
		return err
	}
	return nil
}
