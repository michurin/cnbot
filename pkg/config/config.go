package config

import (
//	"fmt"
//	"github.com/spf13/viper"
	"net/url"
	"time"
)

type BotConfiguration struct {
	Name string
	Token string
	PollingInterval int
	Timeout time.Duration
	Proxy *url.URL
	ControlServer string
}

func GetConfiguration() []BotConfiguration {
	proxy, err := url.Parse("socks5://localhost:8090")
	if err != nil {
		panic(err)  // TODO
	}
	return []BotConfiguration{{
		Name: "main",
		Token: "226015286:AAHzO9VmmKi-_uwZkmbA0DVn03DOpdYYsg4",
		PollingInterval: 20,
		Timeout: 30 * time.Second,
		Proxy: proxy,
	}}
}

/*
type configuration struct {
	Test int `yaml:"test"`
	Network struct {
		Proxy string `yaml:"proxy"`
		PollingInterval int
		LogLevel int
	} `yaml:"network"`
}

func GetConfiguration() int {
	viper.Debug()
	viper.AddConfigPath(".")
	viper.SetConfigName("cnbot")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	cfg := &configuration{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Println(cfg)
	fmt.Println(viper.GetInt("test"))
	fmt.Println(viper.GetStringMap("bots"))
	viper.Debug()
	return 1
}
*/


