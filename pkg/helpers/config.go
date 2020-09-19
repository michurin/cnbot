package helpers

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Bots []BotConfig `json:"bots"`
}

type BotConfig struct {
	Token        string `json:"token"`
	AllowedUsers []int  `json:"allowed_users"`
	Script       string `json:"script"`
	WorkingDir   string `json:"working_dir"`
}

func ReadConfig() Config {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	cfg := new(Config)
	err = json.Unmarshal(data, cfg)
	if err != nil {
		panic(err)
	}
	return *cfg
}
