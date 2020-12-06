package bot

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

var /* const */ safeChars = regexp.MustCompile("[a-zA-Z0-9._-]+")

func argsSafe(s string) []string {
	return safeChars.FindAllString(strings.ToLower(s), -1)
}

func envList(m tg.Message, target int64, server string) []string {
	e := []string{
		"BOT_NAME",
		m.BotName,
		"BOT_FROM",
		hps.Itoa(target),
		"BOT_FROM_FIRSTNAME",
		m.FromFirstName,
		"BOT_CHAT",
		hps.Itoa(m.ChatID),
		"BOT_TEXT",
		m.Text,
	}
	if server != "" {
		e = append(e,
			"BOT_SERVER",
			server,
		)
	}
	if m.SideType != "" {
		e = append(e,
			"BOT_SIDE_TYPE",
			m.SideType,
			"BOT_SIDE_ID",
			hps.Itoa(m.SideID),
			"BOT_SIDE_NAME",
			m.SideName,
		)
	}
	return e
}

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
		argsSafe(m.Text),
		hps.Env(envList(m, target, bot.BindAddress)...),
		bot.WorkingDir)
	if err != nil {
		hps.Log(ctx, err)
		return
	}
	hps.Log(ctx, "script output:", stdout)
	err = SmartSend(ctx, bot.Token, target, m.UpdateMessageID, stdout)
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
