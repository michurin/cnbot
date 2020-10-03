package bot

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	hps "github.com/michurin/cnbot/pkg/helpers"
)

type Handler struct {
	BotName      string
	Token        string
	AllowedUsers map[int]struct{}
}

func (h *Handler) pathDecode(ctx context.Context, path string) (destUser int, err error) {
	urlParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(urlParts) != 1 {
		err = errors.New("invalid path")
		hps.Log(ctx, path, err)
		return
	}
	user := urlParts[len(urlParts)-1]
	destUser, err = strconv.Atoi(user)
	if err != nil {
		hps.Log(ctx, user, err)
		return
	}
	if _, ok := h.AllowedUsers[destUser]; !ok {
		err = errors.New("user is not allowed")
		hps.Log(ctx, destUser, err)
		return
	}
	return
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := hps.Label(r.Context(), hps.RandLabel(), h.BotName)
	hps.Log(ctx, r.Method, r.URL.String())
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	destUser, err := h.pathDecode(ctx, r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		hps.Log(ctx, r.URL.String(), err)
		return
	}
	ctx = hps.Label(ctx, destUser)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		hps.Log(ctx, body, err)
		return
	}
	err = SmartSend(ctx, h.Token, destUser, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		hps.Log(ctx, body, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	hps.Log(ctx, http.StatusOK)
}

func RunHTTPServer(ctx context.Context, addr string, writeTimeout time.Duration, readTimeout time.Duration, handler http.Handler) {
	ctx = hps.Label(ctx, addr)
	s := http.Server{
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		ErrorLog:     log.New(os.Stdout, "http", log.LstdFlags|log.Llongfile|log.Lmsgprefix), // TODO establish wrapper for helpers/log.go
		Addr:         addr,
		Handler:      handler,
	}
	go func() { // what if we shutdown before listen?
		<-ctx.Done()
		hps.Log(ctx, "Server is going to shutdown")
		dCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := s.Shutdown(dCtx)
		if err != nil {
			hps.Log(ctx, err)
		}
	}()
	hps.Log(ctx, "Server is starting on", s.Addr, "with timeouts", s.ReadTimeout, s.WriteTimeout)
	hps.Log(ctx, s.ListenAndServe())
	hps.Log(ctx, "Server finished")
}
