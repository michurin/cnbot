package bot

import (
	"context"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

func buildRequest(destUser, callbackMessageID int64, stdout []byte) (req *tg.Request, err error) {
	imgExt := hps.ImageType(stdout)
	if imgExt != "" {
		req, err = tg.EncodeSendPhoto(destUser, imgExt, stdout)
		return
	}
	ignore, msg, isMarkdown, forUpdate, markup, err := hps.MessageType(stdout)
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
	if callbackID != "" {
		// Slightly hackish. We have to make answerCallbackQuery call
		// TODO: It is the simples kind of answer. Text can be added.
		req, err := tg.EncodeAnswerCallbackQuery(callbackID)
		if err != nil {
			hps.Log(ctx, err)
			return err
		}
		body, err := hps.Do(ctx, tg.Encode(token, req))
		if err != nil {
			hps.Log(ctx, err)
			return err
		}
		err = tg.DecodeSendMessage(body)
		if err != nil {
			hps.Log(ctx, err)
			return err
		}
	}
	req, err := buildRequest(destUser, callbackMessageID, stdout)
	if err != nil {
		hps.Log(ctx, err)
		return err
	}
	if req == nil { // message to be ignored
		return nil
	}
	body, err := hps.Do(ctx, tg.Encode(token, req))
	if err != nil {
		hps.Log(ctx, err)
		return err
	}
	err = tg.DecodeSendMessage(body)
	if err != nil {
		hps.Log(ctx, err)
		return err
	}
	return nil
}
