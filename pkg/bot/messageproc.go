package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

func process(ctx context.Context, botMap map[string]hps.BotConfig, m tg.Message) {
	target := m.ChatID
	ctx = hps.Label(ctx, hps.RandLabel(), m.BotName, target)
	hps.Log(ctx, "Message", m.Text)
	bot, ok := botMap[m.BotName]
	if !ok {
		hps.Log(ctx, fmt.Errorf("bot `%s` is not known", m.BotName))
		return
	}
	if _, ok := bot.AllowedUsers[target]; !ok {
		hps.Log(ctx, fmt.Errorf("user %d is not allowed", target))
		return
	}
	stdout, err := hps.Exec(
		ctx,
		bot.ScriptTermTimeout,
		bot.ScriptKillTimeout,
		bot.ScriptWaitTimeout,
		bot.Script,
		strings.Fields(strings.ToLower(m.Text)), // TODO config
		hps.Env(
			"BOT_NAME",
			m.BotName,
			"BOT_FROM",
			strconv.Itoa(target),
			"BOT_FROM_FIRSTNAME",
			m.FromFirstName,
			"BOT_CHAT",
			strconv.Itoa(m.ChatID),
			"BOT_SERVER",
			bot.BindAddress),
		bot.WorkingDir)
	if err != nil {
		hps.Log(ctx, err)
		return
	}
	hps.Log(ctx, "script output:", stdout)
	err = SmartSend(ctx, bot.Token, target, stdout)
	if err != nil {
		hps.Log(ctx, err)
		return
	}
}

func MessageProcessor(ctx context.Context, msgQueue <-chan tg.Message, botMap map[string]hps.BotConfig) {
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
