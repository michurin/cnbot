package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type baseResponse struct {
	Ok          bool            `json:"ok"`
	Description *string         `json:"description"`
	ErrorCode   *int            `json:"error_code"`
	Result      json.RawMessage `json:"result"`
}

func extractResult(body []byte) (json.RawMessage, error) {
	x := baseResponse{}
	err := json.Unmarshal(body, &x)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !x.Ok {
		return nil, errors.New(errorMessage(x))
	}
	return x.Result, nil
}

func errorMessage(x baseResponse) string {
	message := []string(nil)
	if x.ErrorCode != nil {
		message = append(message, fmt.Sprintf("[error_code=%d]", *x.ErrorCode))
	}
	if x.Description != nil {
		message = append(message, *x.Description)
	}
	return strings.Join(message, " ")
}
