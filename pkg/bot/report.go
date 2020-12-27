package bot

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

const version = "2.0.0"

var /* const */ Build = "noBuildInfo" // go build -ldflags "-X github.com/michurin/cnbot/pkg/bot.Build=`date +%F`-`git rev-parse --short HEAD`" ./cmd/...

func serverConfigurationToString(addr string, w, r time.Duration) string {
	if addr == "" {
		return " server is not configured"
	}
	return fmt.Sprintf(`
      - address: %q
      - timeouts: %v, %v (w/r)`,
		addr, w, r)
}

func botInfo(ctx context.Context, token string) string {
	out, err := hps.Do(ctx, tg.Encode(token, tg.EncodeGetMe()))
	if err != nil {
		hps.Log(ctx, err)
		return " ERROR"
	}
	botID, botName, firstName, canJoinGrp, canReadAllGrpMsg, supportInline, err := tg.DecodeGetMe(out)
	if err != nil {
		hps.Log(ctx, err)
		return " ERROR"
	}
	return fmt.Sprintf(`
    - bot id: %d
    - bot name: %q
    - first name: %q
    - can join grp: %v
    - can read all grp msgs: %v
    - support inline: %v`, botID, botName, firstName, canJoinGrp, canReadAllGrpMsg, supportInline)
}

func botWebHook(ctx context.Context, token string) string {
	out, err := hps.Do(ctx, tg.Encode(token, tg.EncodeGetWebhookInfo()))
	if err != nil {
		hps.Log(ctx, err)
		return "ERROR"
	}
	u, err := tg.DecodeGetWebhookInfo(out)
	if err != nil {
		hps.Log(ctx, err)
		return "ERROR"
	}
	if u == "" {
		return "empty (it's ok)"
	}
	return fmt.Sprintf("%q (NOT OK!)", u)
}

func BotsReport(rootCtx context.Context, cfgs map[string]hps.BotConfig) string {
	nicks := make([]string, len(cfgs))
	i := 0
	for nick := range cfgs {
		nicks[i] = nick
		i++
	}
	sort.Strings(nicks)
	reports := make([]string, len(nicks))
	for i, nick := range nicks {
		ctx := hps.Label(rootCtx, nick)
		c := cfgs[nick]
		reports[i] = fmt.Sprintf(`- version: %s-%s
- go version: %s / %s / %s
- nickname: %q
  - bot info:%s
    - web hook: %s
  - configuration:
    - allowed users: %s
    - script:
      - script: %q
      - working dir: %q
      - timeouts: %v, %v, %v (term/kill/wait)
    - server:%s`,
			version,
			Build,
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH,
			nick,
			botInfo(ctx, c.Token),
			botWebHook(ctx, c.Token),
			c.Access.String(),
			c.Script,
			c.WorkingDir,
			c.ScriptTermTimeout,
			c.ScriptKillTimeout,
			c.ScriptWaitTimeout,
			serverConfigurationToString(c.BindAddress, c.WriteTimeout, c.ReadTimeout))
	}
	return strings.Join(reports, "\n")
}
