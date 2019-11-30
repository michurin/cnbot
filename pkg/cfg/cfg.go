package cfg

import (
	"fmt"
	"os"

	"github.com/michurin/cnbot/pkg/interfaces"
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
)

type BotConfig struct {
	Name   string
	Token  string
	Script string
}

func (b BotConfig) String() string {
	return fmt.Sprintf(
		"\tName: %s\n\tToken: %s\n\tScript: %s",
		b.Name,
		hideToken(b.Token),
		b.Script)
}

type cfgSection struct {
	Token  *string `ini:"token"`
	Script string  `ini:"script"`
}

func Read(fileName string, logger interfaces.Logger) ([]BotConfig, error) {
	f, err := ini.Load(fileName)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	f.ValueMapper = os.ExpandEnv
	sections := f.SectionStrings()
	configs := []BotConfig(nil)
	for _, sectionName := range sections {
		c := new(cfgSection)
		err := f.Section(sectionName).MapTo(c)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if c.Token == nil {
			logger.Log("Left section " + sectionName + ": no token")
			continue
		}
		configs = append(configs, BotConfig{
			Name:   sectionName,
			Token:  *c.Token,
			Script: pathToScript(fileName, c.Script),
		})
	}
	return configs, nil
}

func hideToken(s string) string {
	l := len(s)
	if l < 16 {
		return "[hm... too short]"
	}
	return s[:4] + "..." + s[l-4:]
}
