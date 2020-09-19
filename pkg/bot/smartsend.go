package bot

import (
	"context"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

func buildRequest(destUser int, stdout []byte) (req *tg.Request, err error) {
	imgExt := hps.ImageType(stdout)
	if imgExt != "" {
		req, err = tg.EncodeSendPhoto(destUser, imgExt, stdout)
		return
	}
	ignore, msg, isMarkdown, err := hps.MessageType(stdout)
	if err != nil {
		return
	}
	if ignore {
		return nil, nil
	}
	req, err = tg.EncodeSendMessage(destUser, msg, isMarkdown)
	if err != nil {
		return
	}
	return
}

func SmartSend(ctx context.Context, token string, destUser int, stdout []byte) error {
	req, err := buildRequest(destUser, stdout)
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
