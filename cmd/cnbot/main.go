package main

// bot name @M_78c9409716d3f0bdfd6d_bot

import (
	"github.com/michurin/cnbot/pkg/cnbot"
	"github.com/michurin/cnbot/pkg/cnbot/config"
)

func main() {
	proxyServer := config.GetConfiguration()
	cnbot.PollingLoop(proxyServer)
}