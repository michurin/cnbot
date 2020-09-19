package bot

import (
	"context"
	"sync"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

func Run(ctx context.Context) {
	hps.Log(ctx, "Bot is starting...") // TODO log bot version

	var wg sync.WaitGroup

	msgQueue := make(chan tg.Message)

	configBots, err := hps.ReadConfig()
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

	for botName, bot := range bots {
		wg.Add(1)
		go func(n string, b hps.BotConfig) {
			defer wg.Done()
			Poller(ctx, n, b, msgQueue)
		}(botName, bot)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		RunHTTPServer(ctx, bots)
	}()

	MessageProcessor(ctx, msgQueue, bots)

	wg.Wait()

	hps.Log(ctx, "Bot has been stopped")
}
