package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	BindAddress       string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
}

type config struct {
	Bots map[string]botConfig `json:"bots"`
}

type botConfig struct {
	Token          string   `json:"token"`
	AllowedUsers   []int    `json:"allowed_users"`
	Script         string   `json:"script"`
	WorkingDir     string   `json:"working_dir"`
	TermTimeout    *float64 `json:"term_timeout"`
	KillTimeout    *float64 `json:"kill_timeout"`
	WaitTimeout    *float64 `json:"wait_timeout"`
	BindingAddress string   `json:"bind_address"`
	ReadingTimeout *float64 `json:"read_timeout"`
	WritingTimeout *float64 `json:"write_timeout"`
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

func ReadConfig(configFile string) (map[string]BotConfig, error) {
	if !filepath.IsAbs(configFile) {
		return nil, errors.New("path to config file must be absolute")
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	cfg := new(config)
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	baseDir := filepath.Dir(configFile)
	botCfg := map[string]BotConfig{}
	for nick, b := range cfg.Bots {
		au, err := allowedUsers(b.AllowedUsers)
		if err != nil {
			return nil, err
		}
		script := toAbsPath(baseDir, b.Script)
		stat, err := os.Stat(script)
		if err != nil {
			return nil, err
		}
		if !stat.Mode().IsRegular() {
			return nil, fmt.Errorf("script file %s is not regular", script)
		}
		if stat.Mode()&0111 == 0 { // slightly weird. just mimic exec.findExecutable
			return nil, fmt.Errorf("script file %s is not permited", script)
		}
		pwd := toAbsPath(baseDir, b.WorkingDir)
		stat, err = os.Stat(pwd)
		if err != nil {
			return nil, err
		}
		if !stat.Mode().IsDir() {
			return nil, fmt.Errorf("working dir %s is not a dirrectory", pwd)
		}
		botCfg[nick] = BotConfig{
			Token:             b.Token, // TODO check not empty? some format?
			AllowedUsers:      au,
			Script:            script,
			WorkingDir:        pwd,
			ScriptTermTimeout: defaultDuration(b.TermTimeout, 10*time.Second),
			ScriptKillTimeout: defaultDuration(b.KillTimeout, time.Second),
			ScriptWaitTimeout: defaultDuration(b.WaitTimeout, time.Second),
			BindAddress:       b.BindingAddress,
			ReadTimeout:       defaultDuration(b.ReadingTimeout, 10*time.Second),
			WriteTimeout:      defaultDuration(b.WritingTimeout, 10*time.Second),
		}
	}
	return botCfg, nil
}

func defaultDuration(dur *float64, def time.Duration) time.Duration {
	if dur == nil {
		return def
	}
	return time.Duration(*dur * float64(time.Second))
}
