package cfg

import (
	"flag"

	"github.com/michurin/cnbot/pkg/log"
	"github.com/pelletier/go-toml"
)

type botConfig struct {
	ArgsTrim        *bool    `toml:"args_trim"`
	ArgsLowerCase   *bool    `toml:"args_lower_case"`
	ArgsSplit       *bool    `toml:"args_split"`
	Command         string   `toml:"command"`
	Concurrent      int      `toml:"concurrent"`
	Cwd             string   `toml:"cwd"`
	EnvForce        []string `toml:"env_force"`
	EnvPass         []string `toml:"env_pass"`
	PollingInterval int      `toml:"polling_interval"`
	Port            int64    `toml:"port"`
	ReplayToUser    bool     `toml:"replay_to_user"`
	Timeout         int64    `toml:"timeout"`
	Token           string   `toml:"token"`
	WhiteList       []int64  `toml:"whitelist"`
}

func ReadConfig(log *log.Logger) map[string]botConfig {
	var config_file string
	flag.StringVar(&config_file, "C", "config.toml", "Path to TOML configuration file")
	flag.Parse()
	log.Debugf("Read configuration from %s", config_file)
	tree, err := toml.LoadFile(config_file)
	if err != nil {
		log.Fatal(err)
	}
	cfg := map[string]botConfig{}
	for _, k := range tree.Keys() {
		log.Debugf("Reading section %s", k)
		section := tree.Get(k).(*toml.Tree)
		item := botConfig{}
		err = section.Unmarshal(&item)
		if err != nil {
			log.Fatal(err)
		}
		if item.Token == "" {
			log.Fatalf("You have to specify token for bot %s", k)
		}
		if item.Command == "" {
			log.Fatalf("You have to specify command for bot %s", k)
		}
		if item.Cwd == "" {
			item.Cwd = "."
			log.Warnf("Use cwd='%s' for bot %s", item.Cwd, k)
		}
		if len(item.WhiteList) == 0 {
			log.Fatalf("Whitelist for bot %s is empty", k)
		}
		if len(item.EnvPass)+len(item.EnvForce) == 0 { // TODO: check
			log.Warnf("No envs to pass or force for bot %s", k)
		}
		if item.Timeout == 0 {
			item.Timeout = 5
		}
		if item.Timeout < 0 || item.Timeout > 600 {
			log.Fatalf("Invalid timeout for bot %s", k)
		}
		if item.Concurrent == 0 {
			item.Concurrent = 2
		}
		if item.Concurrent < 0 || item.Concurrent > 100 {
			log.Fatalf("Invalid concurrent for bot %s", k)
		}
		if item.Port == 0 {
			item.Port = 3003
		}
		if item.Port < 1 {
			log.Fatalf("Invalid port for bot %s", k)
		}
		if item.PollingInterval == 0 {
			item.PollingInterval = 50
			log.Warnf("Force default polling interval %d", item.PollingInterval)
		}
		if item.PollingInterval < 10 {
			item.PollingInterval = 10
			log.Warnf("Force polling interval %d", item.PollingInterval)
		}
		if item.PollingInterval > 60 {
			item.PollingInterval = 60
			log.Warnf("Force polling interval %d", item.PollingInterval)
		}
		t := true
		if item.ArgsTrim == nil {
			item.ArgsTrim = &t
		}
		if item.ArgsLowerCase == nil {
			item.ArgsLowerCase = &t
		}
		if item.ArgsSplit == nil {
			item.ArgsSplit = &t
		}
		cfg[k] = item
	}
	log.Debugf("Config: %#v", cfg)
	return cfg
}
