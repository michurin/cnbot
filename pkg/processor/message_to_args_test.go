package processor

import (
	"testing"

	"github.com/michurin/cnbot/pkg/receiver"
)

func slicesNE(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return true
	}
	if len(a) != len(b) {
		return true
	}
	for i := range a {
		if a[i] != b[i] {
			return true
		}
	}
	return false
}

func s(a string) *string {
	return &a
}

func TestMessageToArgs(t *testing.T) {
	for _, c := range []struct {
		Title      string
		MA         MessageToArgs
		Update     receiver.TUpdateResult
		IsCallBack bool
		Args       []string
		Error      error
	}{
		{
			"Empty",
			MessageToArgs{false, false, false},
			receiver.TUpdateResult{},
			false,
			nil,
			errorNoBody,
		},
		{
			"Text: Keep orig",
			MessageToArgs{false, false, false},
			receiver.TUpdateResult{Message: &receiver.TUpdateMessage{Text: s(" AB C ")}},
			false,
			[]string{" AB C "},
			nil,
		},
		{
			"Text: Trim",
			MessageToArgs{true, false, false},
			receiver.TUpdateResult{Message: &receiver.TUpdateMessage{Text: s(" AB C ")}},
			false,
			[]string{"AB C"},
			nil,
		},
		{
			"Text: Trim LowerCase Split",
			MessageToArgs{true, true, true},
			receiver.TUpdateResult{Message: &receiver.TUpdateMessage{Text: s(" AB C ")}},
			false,
			[]string{"ab", "c"},
			nil,
		},
		{
			"Text: LowerCase Split",
			MessageToArgs{false, true, true},
			receiver.TUpdateResult{Message: &receiver.TUpdateMessage{Text: s(" AB C ")}},
			false,
			[]string{"ab", "c"},
			nil,
		},
		{
			"Callback: nil (invalid message)",
			MessageToArgs{false, true, true},
			receiver.TUpdateResult{CallbackQuery: &receiver.TUpdateCallbackQuery{}},
			false,
			nil,
			errorNoBody,
		},
		{
			"Callback: Data",
			MessageToArgs{false, true, true},
			receiver.TUpdateResult{CallbackQuery: &receiver.TUpdateCallbackQuery{Data: s(" A B C ")}},
			true,
			[]string{"callback_data: A B C "},
			nil,
		},
	} {
		t.Run(c.Title, func(t *testing.T) {
			isCallBack, args, err := c.MA.MessageToArgs(c.Update)
			if isCallBack != c.IsCallBack {
				t.Errorf("Error in '%s: Flag: fact=%v expected=%v'", c.Title, isCallBack, c.IsCallBack)
			}
			if err != c.Error {
				t.Errorf("Error in '%s': Error: fact=%#v expected=%#v", c.Title, err, c.Error)
			}
			if slicesNE(args, c.Args) {
				t.Errorf("Error in '%s': Args: fact=%#v expected=%#v", c.Title, args, c.Args)
			}
		})
	}
}
