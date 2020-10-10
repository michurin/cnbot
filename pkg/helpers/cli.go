package helpers

import (
	"flag"
)

func CommandLine() (configFile string, botInfo bool, err error) {
	configFilePtr := flag.String("c", "", "Configuration file in JSON format")
	botInfoPtr := flag.Bool("i", false, "Dump bot info and exit")
	flag.Parse()
	configFile = *configFilePtr
	botInfo = botInfoPtr != nil && *botInfoPtr
	return
}