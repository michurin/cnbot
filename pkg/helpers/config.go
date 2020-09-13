package helpers

import (
	"encoding/json"
	"io/ioutil"
)

type BotConfig struct {
	Token  string
	Script string
}

func Config() []BotConfig {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	cfg := []BotConfig(nil)
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
