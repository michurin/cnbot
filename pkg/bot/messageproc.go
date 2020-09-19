package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

func process(ctx context.Context, botMap map[string]Bot, m tg.Message) {
	fromStr := strconv.Itoa(m.FromID)
	ctx = hps.Label(ctx, m.BotName, fromStr, hps.RandLabel())
	hps.Log(ctx, "Message", m.FromID, m.BotName, m.Text)
	bot, ok := botMap[m.BotName]
	if !ok {
		hps.Log(ctx, fmt.Errorf("bot `%s` is not known", m.BotName))
		return
	}
	if _, ok := bot.AllowedUsers[m.FromID]; !ok {
		hps.Log(ctx, fmt.Errorf("user %d is not allowed", m.FromID))
		return
	}
	stdout, stderr, err := hps.Exec(
		ctx,
		time.Second*10, // TODO config
		time.Second,    // TODO config
		time.Second,    // TODO config
		bot.Script,
		strings.Fields(strings.ToLower(m.Text)), // TODO config
		hps.Env("BOT_NAME", m.BotName, "BOT_FROM", fromStr),
		bot.WorkingDir)
	if len(stderr) > 0 {
		hps.Log(ctx, stderr)
	}
	if err != nil {
		hps.Log(ctx, err)
		return
	}
	hps.Log(ctx, "script output:", stdout)
	err = SmartSend(ctx, bot.Token, m.FromID, stdout)
	if err != nil {
		hps.Log(ctx, err)
		return
	}
}

func MessageProcessor(ctx context.Context, msgQueue <-chan tg.Message, botMap map[string]Bot) {
	for {
		select {
		case <-ctx.Done():
			hps.Log(ctx, "Queue listener exited due to context cancellation")
			return
		case m := <-msgQueue:
			process(ctx, botMap, m)
		}
	}
}
