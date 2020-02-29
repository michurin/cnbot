package apirequest

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func EncodeJSON(params interface{}) (Request, error) {
	var err error
	var body []byte
	if params != nil {
		body, err = json.Marshal(params)
		if err != nil {
			return Request{}, errors.WithStack(err)
		}
	}
	return Request{
		Method: "POST",
		MIME:   "application/json",
		Body:   body,
	}, nil
}
