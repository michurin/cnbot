package calltgapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/michurin/cnbot/pkg/log"
)

func url(token string, method string) string {
	return "https://api.telegram.org/bot" + token + "/" + method
}

func PostBytes(log *log.Logger, token string, method string, data []byte, mime string, resp interface{}) error {
	//log.Debugf("Raw send: %v", data)
	response, err := http.Post(url(token, method), mime, bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	//log.Debugf("Raw recv: %v", body)
	err = json.Unmarshal(body, resp)
	return err
}

func PostStruct(log *log.Logger, token string, method string, req interface{}, resp interface{}) error { // TODO: called only in one place?
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return PostBytes(log, token, method, body, "application/json", resp)
}
