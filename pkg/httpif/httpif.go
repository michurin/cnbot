package httpif

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/michurin/cnbot/pkg/log"
	"github.com/michurin/cnbot/pkg/prepareoutgoing"
	"github.com/michurin/cnbot/pkg/sender"
)

func valueToMap(q url.Values) (r map[string]string) {
	for k, v := range q {
		if r == nil {
			r = map[string]string{}
		}
		r[k] = v[0]
	}
	return
}

type handler struct {
	outQueue chan sender.OutgoingData
	log      *log.Logger
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chatIdString := r.URL.Path
	chatId, err := strconv.ParseInt(chatIdString[1:], 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.log.Info(chatId)
	body, err := ioutil.ReadAll(r.Body) // The ServeHTTP Handler does not need to close Body.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	k := r.URL.Query()
	resp := make(chan []byte, 1)
	q, err := prepareoutgoing.PrepareOutgoing(h.log, body, chatId, valueToMap(k), resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	replay_data := []byte("nodata")
	if q.MessageType != "" {
		h.outQueue <- q
		replay_data = <-resp
	}
	w.Write(replay_data)
}

func HttpIf(log *log.Logger, port int, oq chan sender.OutgoingData) {
	s := &http.Server{
		Addr:           ":" + strconv.Itoa(port),
		Handler:        handler{oq, log},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 12,
	}
	log.Fatal(s.ListenAndServe())
}
