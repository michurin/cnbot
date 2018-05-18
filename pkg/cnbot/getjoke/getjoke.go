package getjoke

import (
	"github.com/michurin/cnbot/pkg/cnbot/httpcall"
	"net/http"
	"github.com/michurin/cnbot/pkg/perror"
	"encoding/json"
	"html"
)

const API_CALL_ICNDB = "http://api.icndb.com/jokes/random?limitTo=[nerdy]"

type resultStruct struct {  // не экспортится
	Type string
	Value struct {
		Id int
		Joke string
		Categories []string
	}
}

func GetJoke() (string, error) {
	var resultStruct resultStruct
	body, err := httpcall.HTTPCall(http.DefaultClient, http.MethodGet, API_CALL_ICNDB, nil)  // TODO клиент должен быть вынесен в аргументы
	if err != nil {
		return "", perror.NewErrorString("HTTP error: %s", err)
	}
	err = json.Unmarshal(body, &resultStruct)
	if err != nil {
		return "", perror.NewErrorString("JSON parse error: %s", err)
	}
	//fmt.Printf("%#v\n", resultStruct)
	if resultStruct.Type != "success" {
		return "", perror.NewErrorString("Infalid type '%s'", resultStruct.Type)
	}
	return html.UnescapeString(resultStruct.Value.Joke), nil
}