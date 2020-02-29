package cfg

import (
	"fmt"
	"strings"
)

type AppConfig struct {
	CheckMode bool // TODO
	BindAddress string // TODO
	Bots []BotConfig
}

type BotConfig struct {
	Name         string
	Token        string
	Script       string
	Env          []string
	AllowedUsers []int
}

func (b BotConfig) String() string {
	return fmt.Sprintf(
		"\tName: %s\n\tToken: %s\n\tScript: %s\n\tEnvs: %s",
		b.Name,
		hideToken(b.Token),
		b.Script,
		strings.Join(b.Env, ":"))
}

func hideToken(s string) string {
	l := len(s)
	if l < 16 {
		return "[hm... too short]"
	}
	return s[:4] + "..." + s[l-4:]
}
