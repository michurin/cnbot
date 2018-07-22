package main

import "github.com/michurin/cnbot/pkg/bot"

var (
	Version   = "n/a"
	BuildRev  = "n/a"
	BuildDate = "n/a"
)

func main() {
	bot.Run(Version, BuildRev, BuildDate)
}
