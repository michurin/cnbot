package cfg

import (
	"fmt"
	"os"
	"strings"

	"github.com/michurin/cnbot/pkg/interfaces"
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
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

type cfgSection struct {
	Token    *string  `ini:"token"`
	Script   string   `ini:"script"`
	ForceEnv []string `ini:"force_env,omitempty,allowshadow"`
	PassEnv  []string `ini:"pass_env,omitempty,allowshadow"`
	Allowed  []int    `ini:"allowed,omitempty,allowshadow"`
}

func Read(fileName string, logger interfaces.Logger) (AppConfig, error) {
	f, err := ini.ShadowLoad(fileName)
	if err != nil {
		return AppConfig{}, errors.WithStack(err)
	}
	f.ValueMapper = os.ExpandEnv
	sections := f.SectionStrings()
	configs := []BotConfig(nil)
	for _, sectionName := range sections {
		c := new(cfgSection)
		err := f.Section(sectionName).MapTo(c)
		if err != nil {
			return AppConfig{}, errors.WithStack(err)
		}
		if c.Token == nil {
			logger.Log("Left section " + sectionName + ": no token")
			continue
		}
		configs = append(configs, BotConfig{
			Name:         sectionName,
			Token:        *c.Token,
			AllowedUsers: c.Allowed,
			Script:       pathToScript(fileName, c.Script),
			Env:          prepareEnv(c.PassEnv, c.ForceEnv),
		})
	}
	return AppConfig{
		Bots: configs,
	}, nil
}

func hideToken(s string) string {
	l := len(s)
	if l < 16 {
		return "[hm... too short]"
	}
	return s[:4] + "..." + s[l-4:]
}
