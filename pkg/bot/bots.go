package bot

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	hps "github.com/michurin/cnbot/pkg/helpers"
	"github.com/michurin/cnbot/pkg/tg"
)

func allowedUsersToString(a map[int]struct{}) string {
	if len(a) == 0 {
		return "empty (nobody can use this bot)"
	}
	v := make([]int, len(a))
	i := 0
	for k := range a {
		v[i] = k
		i++
	}
	sort.Ints(v)
	w := make([]string, len(v))
	sep := ""
	if len(v) < 10 {
		for i, u := range v {
			w[i] = " " + strconv.Itoa(u)
		}
		sep = ","
	} else {
		for i, u := range v {
			w[i] = "\n      - " + strconv.Itoa(u)
		}
	}
	return strings.Join(w, sep)
}

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
	botID, botName, err := tg.DecodeGetMe(out)
	if err != nil {
		hps.Log(ctx, err)
		return " ERROR"
	}
	return fmt.Sprintf(`
    - bot id: %d
    - bot name: %q`, botID, botName)
}

func BotsReport(rootCtx context.Context, cfgs map[string]hps.BotConfig) (string, error) {
	reports := make([]string, len(cfgs))
	i := 0
	for nick, c := range cfgs { // TODO sort by nickname
		ctx := hps.Label(rootCtx, nick)
		reports[i] = fmt.Sprintf(`- nickname: %q
  - bot info:%s
  - configuration:
    - allowed users:%s
    - script:
      - script: %q
      - working dir: %q
      - timeouts: %v, %v, %v (term/kill/wait)
    - server:%s`,
			nick,
			botInfo(ctx, c.Token),
			allowedUsersToString(c.AllowedUsers),
			c.Script,
			c.WorkingDir,
			c.ScriptTermTimeout,
			c.ScriptKillTimeout,
			c.ScriptWaitTimeout,
			serverConfigurationToString(c.BindAddress, c.WriteTimeout, c.ReadTimeout))
		i++
	}
	return strings.Join(reports, "\n"), nil
}
