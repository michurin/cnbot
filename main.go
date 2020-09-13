package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

func main() {
	//stdout, stderr, exitCode, err := hps.Exec(
	//	context.Background(),
	//	time.Millisecond*500,
	//	time.Millisecond*500,
	//	time.Millisecond*500,
	//	"./script.sh",
	//	nil,
	//	nil,
	//	"")
	//hps.Log(context.Background(), "======")
	//hps.Log(context.Background(), "OUT:", stdout)
	//hps.Log(context.Background(), "ERR:", stderr)
	//hps.Log(context.Background(), "CODE:", exitCode)
	//hps.Log(context.Background(), err)
	//hps.Log(context.Background(), "======")

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
			tg.Poller(ctx, bot, msgQueue)
			wg.Done()
		}(bot)
	}

	botNameToToken := map[string]tg.Bot{ // TODO
		bots[0].Username: bots[0],
	}

MainLoop: // TODO: move to separate func
	for {
		var msg string
		select {
		case <-ctx.Done():
			hps.Log(ctx, "Queue listener exited due to context cancellation")
			break MainLoop
		case m := <-msgQueue:
			hps.Log(ctx, "MESSAGE", m.BotName, m.Text)
			stdout, stderr, exitCode, err := hps.Exec(
				ctx,
				time.Millisecond*1000,
				time.Millisecond*500,
				time.Millisecond*500,
				botNameToToken[m.BotName].Script,
				nil,
				nil,
				"")
			if err == nil {
				msg = fmt.Sprintf("%s [%d]: %s (%s)", botNameToToken[m.BotName].Script, exitCode, stdout, stderr)
			} else {
				msg = err.Error()
				hps.Log(ctx, err)
			}
			req, err := tg.Encode(botNameToToken[m.BotName].Token, tg.EncodeSendMessage(m.FromID, msg))
			if err != nil {
				hps.Log(ctx, err)
				panic(err)
			}
			_, err = hps.Do(ctx, req)
			if err != nil {
				hps.Log(ctx, err)
				panic(err)
			}
		}
	}

	wg.Wait()
}
