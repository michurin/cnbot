package bot

import (
	"os"
	"strconv"

	"github.com/michurin/cnbot/pkg/cfg"
	"github.com/michurin/cnbot/pkg/httpif"
	"github.com/michurin/cnbot/pkg/log"
	"github.com/michurin/cnbot/pkg/processor"
	"github.com/michurin/cnbot/pkg/receiver"
	"github.com/michurin/cnbot/pkg/sender"
)

func Run() {
	log := log.New()
	log.Info("Run")
	config := cfg.ReadConfig(log.WithArea("config"))
	for k, s := range config {
		log.Infof("Going up bot %s", k)
		incomingQueue := make(chan receiver.TUpdateMessage, 1000) // TODO: to config
		outgoingQueue := make(chan sender.OutgoingData, 100)      // TODO: ?
		go receiver.RunPollingLoop(log.WithArea(k+":poller"), s.Token, incomingQueue)
		for p := 0; p < s.Concurrent; p++ {
			go processor.Processor(
				log.WithArea(k+":proc"),
				incomingQueue,
				outgoingQueue,
				s.WhiteList,
				s.Command,
				s.Cwd,
				processor.BuildEnv(os.Environ(), s.EnvPass, append(s.EnvForce, "BOT_SERVER_PORT="+strconv.FormatInt(s.Port, 10))),
				s.Timeout,
			)
		}
		go sender.Sender(log.WithArea(k+":sender"), s.Token, outgoingQueue)
		go httpif.HttpIf(log.WithArea(k+":http"), s.Port, outgoingQueue)
	}
	<-make(chan interface{})
}
