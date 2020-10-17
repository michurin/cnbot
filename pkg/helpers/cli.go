package helpers

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
)

func CommandLine() (configFile string, botInfo bool, err error) {
	configFilePtr := flag.String("c", "", "Configuration file in JSON format")
	botInfoPtr := flag.Bool("i", false, "Dump bot info and exit")
	flag.Parse()
	if configFilePtr == nil {
		err = errors.New("-c is mandatory option")
		return
	}
	configFile = *configFilePtr
	if configFile == "" {
		err = errors.New("configuration file have to be specified (-c)")
		return
	}
	if !filepath.IsAbs(configFile) {
		var executable string
		executable, err = os.Executable()
		if err != nil {
			return
		}
		executable, err = filepath.EvalSymlinks(executable)
		if err != nil {
			return
		}
		configFile = filepath.Join(filepath.Dir(executable), configFile)
	}
	botInfo = botInfoPtr != nil && *botInfoPtr
	return
}
