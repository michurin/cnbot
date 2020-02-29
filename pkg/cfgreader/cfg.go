package cfgreader

import (
	"github.com/michurin/cnbot/pkg/cfg"
	"github.com/michurin/cnbot/pkg/interfaces"
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
	"os"
)

type cfgSection struct {
	Token    *string  `ini:"token"`
	Script   string   `ini:"script"`
	ForceEnv []string `ini:"force_env,omitempty,allowshadow"`
	PassEnv  []string `ini:"pass_env,omitempty,allowshadow"`
	Allowed  []int    `ini:"allowed,omitempty,allowshadow"`
}

func Read(fileName string, logger interfaces.Logger) (cfg.AppConfig, error) {
	f, err := ini.ShadowLoad(fileName)
	if err != nil {
		return cfg.AppConfig{}, errors.WithStack(err)
	}
	f.ValueMapper = os.ExpandEnv
	sections := f.SectionStrings()
	configs := []cfg.BotConfig(nil)
	for _, sectionName := range sections {
		c := new(cfgSection)
		err := f.Section(sectionName).MapTo(c)
		if err != nil {
			return cfg.AppConfig{}, errors.WithStack(err)
		}
		if c.Token == nil {
			logger.Log("Left section " + sectionName + ": no token")
			continue
		}
		configs = append(configs, cfg.BotConfig{
			Name:         sectionName,
			Token:        *c.Token,
			AllowedUsers: c.Allowed,
			Script:       pathToScript(fileName, c.Script),
			Env:          prepareEnv(c.PassEnv, c.ForceEnv),
		})
	}
	return cfg.AppConfig{
		Bots: configs,
	}, nil
}