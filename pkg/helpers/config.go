package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type BotConfig struct {
	Token             string
	AllowedUsers      map[int]struct{}
	Script            string
	WorkingDir        string
	ScriptTermTimeout time.Duration
	ScriptKillTimeout time.Duration
	ScriptWaitTimeout time.Duration
}

type ServerConfig struct {
	BindAddress  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type config struct {
	Bots   []botConfig   `json:"bots"`
	Server *serverConfig `json:"server"`
}

type botConfig struct {
	Token        string `json:"token"`
	AllowedUsers []int  `json:"allowed_users"`
	Script       string `json:"script"`
	WorkingDir   string `json:"working_dir"`
}

type serverConfig struct {
	BindingAddress string `json:"bind_address"`
	ReadingTimeout int    `json:"read_timeout"`
	WritingTimeout int    `json:"write_timeout"`
}

func allowedUsers(uu []int) (map[int]struct{}, error) {
	m := map[int]struct{}{}
	for _, u := range uu {
		if _, ok := m[u]; ok {
			return nil, fmt.Errorf("user %d is allowed twice", u)
		}
		m[u] = struct{}{}
	}
	return m, nil
}

func toAbsPath(baseDir string, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(baseDir, path)
}

func ReadConfig() ([]BotConfig, *ServerConfig, error) {
	configFile := flag.String("c", "config.json", "Configuration file in JSON format")
	flag.Parse()
	if configFile == nil {
		return nil, nil, errors.New("can not receive config file name") // it's impossible
	}
	if !filepath.IsAbs(*configFile) {
		return nil, nil, errors.New("path to config file must be absolute")
	}
	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return nil, nil, err
	}
	cfg := new(config)
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, nil, err
	}
	baseDir := filepath.Dir(*configFile)
	botCfg := make([]BotConfig, len(cfg.Bots))
	for i, b := range cfg.Bots {
		au, err := allowedUsers(b.AllowedUsers)
		if err != nil {
			return nil, nil, err
		}
		botCfg[i] = BotConfig{
			Token:             b.Token, // TODO check not empty? some format?
			AllowedUsers:      au,
			Script:            toAbsPath(baseDir, b.Script),
			WorkingDir:        toAbsPath(baseDir, b.WorkingDir),
			ScriptTermTimeout: 10 * time.Second,
			ScriptKillTimeout: time.Second,
			ScriptWaitTimeout: time.Second,
		}
	}
	if cfg.Server == nil {
		return botCfg, nil, nil
	}
	serverCfg := &ServerConfig{ // TODO check? defaults?
		BindAddress:  cfg.Server.BindingAddress,
		ReadTimeout:  time.Duration(cfg.Server.ReadingTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WritingTimeout) * time.Second,
	}
	return botCfg, serverCfg, nil
}

func allowedUsersToString(uu map[int]struct{}) string {
	v := []int(nil)
	for u := range uu {
		v = append(v, u)
	}
	sort.Ints(v)
	s := make([]string, len(v))
	for i, p := range v {
		s[i] = strconv.Itoa(p)
	}
	return strings.Join(s, ", ")
}

func DumpBotConfig(ctx context.Context, cfg map[string]BotConfig) {
	for n, b := range cfg {
		c := Label(ctx, n)
		Log(c, "Token:", b.Token)
		Log(c, "Allowed users:", allowedUsersToString(b.AllowedUsers))
		Log(c, "Script:", b.Script)
		Log(c, "Working dir:", b.WorkingDir)
	}
}
