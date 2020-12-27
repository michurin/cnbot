package helpers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type BotConfig struct {
	Token             string
	Access            AccessCtl
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
	Bots  map[string]botConfig `yaml:"bots"`
	Alive aliveConfig          `yaml:"alive"`
}

type botConfig struct {
	Token           string    `yaml:"token"`
	AllowedUsers    yaml.Node `yaml:"allowed_users"`
	DisallowedUsers []int64   `yaml:"disallowed_users"`
	Script          string    `yaml:"script"`
	WorkingDir      string    `yaml:"working_dir"`
	TermTimeout     *float64  `yaml:"term_timeout"`
	KillTimeout     *float64  `yaml:"kill_timeout"`
	WaitTimeout     *float64  `yaml:"wait_timeout"`
	BindingAddress  string    `yaml:"bind_address"`
	ReadingTimeout  *float64  `yaml:"read_timeout"`
	WritingTimeout  *float64  `yaml:"write_timeout"`
}

type aliveConfig struct {
	BindingAddress string `yaml:"bind_address"`
}

func allowedUsers(node yaml.Node) (map[int64]struct{}, error) {
	// it accepts
	// list of integers and list of strings
	if node.Kind == 0 {
		return nil, nil
	}
	if node.Kind != yaml.SequenceNode {
		return nil, errors.New("allowed users list has to be sequence")
	}
	m := map[int64]struct{}{}
	for _, n := range node.Content {
		if n.Kind != yaml.ScalarNode {
			return nil, errors.New("allowed users item has to be scalar")
		}
		u, err := Atoi(substituteEnvironment(n.Value))
		if err != nil {
			return nil, err
		}
		if _, ok := m[u]; ok {
			return nil, fmt.Errorf("user %d is allowed twice", u)
		}
		m[u] = struct{}{}
	}
	return m, nil
}

func disallowedUsers(uu []int64) (map[int64]struct{}, error) {
	m := map[int64]struct{}{}
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

func ReadConfig(configFile string) (map[string]BotConfig, string, error) {
	if !filepath.IsAbs(configFile) {
		return nil, "", errors.New("path to config file must be absolute")
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, "", err
	}
	cfg := new(config)
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, "", err
	}
	baseDir := filepath.Dir(configFile)
	botCfg := map[string]BotConfig{}
	for nick, b := range cfg.Bots {
		au, err := allowedUsers(b.AllowedUsers)
		if err != nil {
			return nil, "", err
		}
		dau, err := disallowedUsers(b.DisallowedUsers)
		if err != nil {
			return nil, "", err
		}
		script := toAbsPath(baseDir, b.Script)
		stat, err := os.Stat(script)
		if err != nil {
			return nil, "", err
		}
		if !stat.Mode().IsRegular() {
			return nil, "", fmt.Errorf("script file %s is not regular", script)
		}
		if stat.Mode()&0111 == 0 { // slightly weird. just mimic exec.findExecutable
			return nil, "", fmt.Errorf("script file %s is not permited", script)
		}
		pwd := toAbsPath(baseDir, b.WorkingDir)
		stat, err = os.Stat(pwd)
		if err != nil {
			return nil, "", err
		}
		if !stat.Mode().IsDir() {
			return nil, "", fmt.Errorf("working dir %s is not a dirrectory", pwd)
		}
		botCfg[nick] = BotConfig{
			Token:             substituteEnvironment(b.Token),
			Access:            NewAccessCtl(au, dau),
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
	return botCfg, cfg.Alive.BindingAddress, nil
}

func substituteEnvironment(t string) string {
	if strings.HasPrefix(t, "${") && strings.HasSuffix(t, "}") {
		if v, ok := os.LookupEnv(t[2 : len(t)-1]); ok {
			return v
		}
	}
	return t
}

func defaultDuration(dur *float64, def time.Duration) time.Duration {
	if dur == nil {
		return def
	}
	return time.Duration(*dur * float64(time.Second))
}
