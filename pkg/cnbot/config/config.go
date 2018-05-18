package config

import (
	"github.com/spf13/viper"
	"fmt"
	"log"
)

func GetConfiguration() string {
	viper.SetDefault("proxy", "")
	viper.AddConfigPath(".")
	viper.SetConfigName("cnbot" )
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	proxyServer := viper.GetString("proxy")
	log.Printf("Use proxy server %s\n", proxyServer)
	return proxyServer
}
