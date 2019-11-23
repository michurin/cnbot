// ONE BIG TODO

package server

import (
	"log"
	"net/http"
)

type HTTPHandler struct{}

func (_ HTTPHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Printf("SERVER %s [NOLOG!]", req.URL) // TODO
}
