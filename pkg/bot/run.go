package bot

import (
	"context"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

func Run(rootCtx context.Context) {
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	hps.Log(ctx, "Bot is starting...") // TODO log bot version

	msgQueue := make(chan tg.Message)

	configBots, configServer, err := hps.ReadConfig()
	if err != nil {
		hps.Log(ctx, err)
		return
	}

	bots, err := Bots(ctx, configBots)
	if err != nil { // canceled context cause err too
		hps.Log(ctx, err)
		return
	}

	hps.DumpBotConfig(ctx, bots)

	done := make(chan struct{}, 1)
	doneCount := 0

	for botName, bot := range bots {
		doneCount++
		go func(n string, b hps.BotConfig) {
			defer func() { done <- struct{}{} }()
			Poller(ctx, n, b, msgQueue)
		}(botName, bot)
	}

	if configServer != nil {
		doneCount++
		go func() {
			defer func() { done <- struct{}{} }()
			RunHTTPServer(ctx, configServer, bots)
		}()
	} else {
		hps.Log(ctx, "Server didn't start. Not configured")
	}

	if len(bots) > 0 {
		doneCount++
		go func() {
			defer func() { done <- struct{}{} }()
			MessageProcessor(ctx, msgQueue, bots)
			done <- struct{}{}
		}()
	}

	if doneCount > 0 {
		<-done // waiting for at least one exit
		doneCount--
		cancel() // cancel all
		for ; doneCount >= 0; doneCount-- {
			<-done
		}
	}

	hps.Log(ctx, "Bot has been stopped")
}
