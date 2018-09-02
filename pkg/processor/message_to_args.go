package processor

import (
	"errors"
	"strings"

	"github.com/michurin/cnbot/pkg/receiver"
)

var errorNoBody = errors.New("Can not find message body")

type MessageToArgs struct { // TODO: flag to remove leading slash
	Trim      bool
	LowerCase bool
	SplitArgs bool
}

func (ma MessageToArgs) MessageToArgs(p receiver.TUpdateResult) (bool, []string, error) {
	if p.Message != nil {
		if p.Message.Text != nil {
			s := *p.Message.Text
			if ma.Trim {
				s = strings.TrimSpace(s)
			}
			if ma.LowerCase {
				s = strings.ToLower(s)
			}
			if ma.SplitArgs {
				return false, strings.Fields(s), nil
			}
			return false, []string{s}, nil
		}
	}
	if p.CallbackQuery != nil {
		if p.CallbackQuery.Data != nil {
			return true, []string{"callback_data:" + *p.CallbackQuery.Data}, nil
		}
	}
	return false, nil, errorNoBody
}
