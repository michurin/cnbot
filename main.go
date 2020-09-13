package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, os.Interrupt)
	go func() {
		sig := <-sigs
		hps.Log(ctx, "Killed by signal", sig)
		cancel()
	}()

	var wg sync.WaitGroup

	msgQueue := make(chan tg.Message)

	bots, err := tg.Bots(ctx, hps.Config())
	if err != nil {
		hps.Log(ctx, err)
		panic(err)
	}

	for _, bot := range bots {
		hps.Log(ctx, "Run poller for bot", bot.Username)
		wg.Add(1)
		go func(bot tg.Bot) {
			defer wg.Done()
			tg.Poller(ctx, bot, msgQueue)
		}(bot)
	}

	botNameToToken := map[string]tg.Bot{ // TODO
		bots[0].Username: bots[0],
	}

	// server

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK\n"))
		if err != nil {
			hps.Log(r.Context(), err)
		}
		hps.Log(r.Context(), r.Method, r.URL.Path)
	})

	wg.Add(1)
	go func() {
		defer wg.Done()
		s := http.Server{
			Addr:    ":9090", // TODO config
			Handler: handler,
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
		err := s.ListenAndServe()
		if err != nil {
			hps.Log(ctx, err)
			return
		}
		hps.Log(ctx, "Server finished")
	}()

	// /server

MainLoop: // TODO: move to separate func
	for {
		select {
		case <-ctx.Done():
			hps.Log(ctx, "Queue listener exited due to context cancellation")
			break MainLoop
		case m := <-msgQueue:
			hps.Log(ctx, "MESSAGE", m.BotName, m.Text)
			stdout, stderr, err := hps.Exec(
				ctx,
				time.Millisecond*1000, // TODO config
				time.Millisecond*500,  // TODO config
				time.Millisecond*500,  // TODO config
				botNameToToken[m.BotName].Script,
				strings.Fields(strings.ToLower(m.Text)), // TODO config
				nil,                                     // TODO config
				"")                                      // TODO config
			if len(stderr) > 0 {
				hps.Log(ctx, stderr)
			}
			if err != nil {
				hps.Log(ctx, err)
				continue
			}
			msg, imgExt, err := tg.DataType(stdout)
			if err != nil {
				hps.Log(ctx, err)
				continue
			}
			var reqX tg.Request
			if imgExt == "" {
				reqX = tg.EncodeSendMessage(m.FromID, msg)
			} else {
				reqX = tg.EncodeSendPhoto(m.FromID, imgExt, stdout)
			}
			req, err := tg.Encode(botNameToToken[m.BotName].Token, reqX)
			if err != nil {
				hps.Log(ctx, err)
				continue
			}
			_, err = hps.Do(ctx, req)
			if err != nil {
				hps.Log(ctx, err)
				continue
			}
		}
	}

	wg.Wait()
}
