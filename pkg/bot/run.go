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

	config := hps.ReadConfig()

	bots, err := Bots(ctx, config.Bots)
	if err != nil {
		hps.Log(ctx, err)
		panic(err)
	}

	for botName, bot := range bots {
		wg.Add(1)
		go func(n string, b Bot) {
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
