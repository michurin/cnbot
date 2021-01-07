package bot

const version = "2.5.4"

var /* const */ Build = "noBuildInfo" // go build -ldflags "-X github.com/michurin/cnbot/pkg/bot.Build=`date +%F`-`git rev-parse --short HEAD`" ./cmd/...

var Version = version

func init() {
	Version = version + "-" + Build
}
