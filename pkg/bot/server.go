package bot

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	hps "github.com/michurin/cnbot/pkg/helpers"
)

type Handler struct {
	BotMap map[string]hps.BotConfig
}

func (h *Handler) pathDecode(ctx context.Context, path string) (destUser int, botName string, bot hps.BotConfig, err error) {
	var ok bool
	urlParts := strings.Split(strings.Trim(path, "/"), "/")
	if len(urlParts) != 3 {
		err = errors.New("invalid path")
		hps.Log(ctx, path, err)
		return
	}
	botName = urlParts[0]
	bot, ok = h.BotMap[botName]
	if !ok {
		err = errors.New("bot does not exists")
		hps.Log(ctx, botName, err)
		return
	}
	act := urlParts[1]
	if act != "to" {
		err = errors.New("invalid action")
		hps.Log(ctx, act, err)
		return
	}
	user := urlParts[2]
	destUser, err = strconv.Atoi(user)
	if err != nil {
		hps.Log(ctx, user, err)
		return
	}
	if _, ok := bot.AllowedUsers[destUser]; !ok {
		err = errors.New("user is not allowed")
		hps.Log(ctx, destUser, err)
		return
	}
	return
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hps.Log(ctx, r.Method, r.URL.String())
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	destUser, botName, bot, err := h.pathDecode(ctx, r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		hps.Log(ctx, r.URL.String(), err)
		return
	}
	ctx = hps.Label(ctx, botName, strconv.Itoa(destUser))
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		hps.Log(ctx, body, err)
		return
	}
	err = SmartSend(ctx, bot.Token, destUser, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		hps.Log(ctx, body, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	hps.Log(ctx, http.StatusOK)
}

func RunHTTPServer(ctx context.Context, cfg *hps.ServerConfig, botMap map[string]hps.BotConfig) {
	s := http.Server{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		ErrorLog:     log.New(os.Stdout, "http", log.LstdFlags|log.Llongfile|log.Lmsgprefix), // TODO establish wrapper for helpers/log.go
		Addr:         cfg.BindAddress,
		Handler:      &Handler{BotMap: botMap},
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return hps.Label(ctx, "["+c.RemoteAddr().String()+"]", hps.RandLabel())
		},
	}
	go func() { // what if we shutdown before listen?
		<-ctx.Done()
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
