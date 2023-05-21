package bot

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/michurin/minlog"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

var /* const */ safeChars = regexp.MustCompile("[a-zA-Z0-9._-]+")

func argsSafe(s string) []string {
	return safeChars.FindAllString(strings.ToLower(s), -1)
}

func msgType(isCallback bool) string {
	if isCallback {
		return "callback"
	}
	return "message"
}

func appendIfNotEmpty(a []string, x ...string) []string {
	for i := 1; i < len(x); i += 2 {
		if x[i] != "" {
			a = append(a, x[i-1], x[i])
		}
	}
	return a
}

func boolToString(a bool) (r string) {
	if a {
		r = "TRUE"
	}
	return
}

func envList(m tg.Message, target int64, server string) []string {
	e := []string{
		"BOT_VERSION", Version,
		"BOT_NAME", m.BotName,
		"BOT_FROM", hps.Itoa(target),
		"BOT_FROM_FIRST_NAME", m.FromFirstName,
		"BOT_CHAT", hps.Itoa(m.ChatID),
		"BOT_TEXT", m.Text,
		"BOT_MESSAGE_TYPE", msgType(m.CallbackID != ""),
	}
	if m.SideType != "" {
		e = append(e,
			"BOT_SIDE_TYPE", m.SideType,
			"BOT_SIDE_ID", hps.Itoa(m.SideID),
			"BOT_SIDE_NAME", m.SideName,
		)
	}
	e = appendIfNotEmpty(e,
		"BOT_SERVER", server,
		"BOT_FROM_LAST_NAME", m.FromLastName,
		"BOT_FROM_USERNAME", m.FromUsername,
		"BOT_FROM_LANGUAGE", m.FromLanguage,
		"BOT_FROM_IS_BOT", boolToString(m.FromIsBot),
		"BOT_LOCATION_LONGITUDE", m.LocationLongitude,
		"BOT_LOCATION_LATITUDE", m.LocationLatitude,
		"BOT_LOCATION_ACCURACY", m.LocationHorizontalAccuracy,
		"BOT_LOCATION_LIVE_PERIOD", m.LocationLivePeriod,
		"BOT_LOCATION_HEADING", m.LocationHeading,
		"BOT_LOCATION_ALERT_RADIUS", m.LocationProximityAlertRadius,
	)
	return e
}

func process(ctx context.Context, botMap map[string]hps.BotConfig, m tg.Message) {
	target := m.ChatID
	ctx = minlog.Label(ctx, m.BotName+":"+hps.AutoLabel()+":"+strconv.FormatInt(target, 10))
	minlog.Log(ctx, "Message", m.Text)
	bot, ok := botMap[m.BotName]
	if !ok {
		minlog.Log(ctx, fmt.Errorf("bot `%s` is not known", m.BotName))
		return
	}
	if !bot.Access.IsAllowed(target) {
		minlog.Log(ctx, fmt.Errorf("user %d is not allowed", target))
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
		minlog.Log(ctx, err)
		return
	}
	minlog.Log(ctx, "script output:", stdout)
	err = SmartSend(ctx, bot.Token, m.CallbackID, target, m.UpdateMessageID, stdout, "")
	if err != nil {
		minlog.Log(ctx, err)
		return
	}
}

func MessageProcessor(ctx context.Context, msgQueue <-chan tg.Message, botMap map[string]hps.BotConfig) {
	for {
		select {
		case <-ctx.Done():
			minlog.Log(ctx, "Queue listener exited due to context cancellation")
			return
		case m := <-msgQueue:
			process(ctx, botMap, m)
		}
	}
}
