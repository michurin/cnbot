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

func Run(version, buildRev, buildDate string) {
	log := log.New()
	log.Infof("Run ver=%s rev=%s date=%s", version, buildRev, buildDate)
	config := cfg.ReadConfig(log.WithArea("config"))
	for k, s := range config {
		log.Infof("Going up bot: %s", k)
		log.Debugf("Bot %s config: %+v", k, s)
		incomingQueue := make(chan receiver.TUpdateResult, 1000) // TODO: to config
		outgoingQueue := make(chan sender.OutgoingData, 100)     // TODO: ?
		go receiver.RunPollingLoop(
			log.WithArea(k+":poller"),
			s.PollingInterval,
			s.Token,
			incomingQueue,
		)
		for p := 0; p < s.Concurrent; p++ {
			go processor.Processor(
				log.WithArea(k+":proc:"+strconv.Itoa(p)),
				incomingQueue,
				outgoingQueue,
				s.WhiteList,
				s.Command,
				s.Cwd,
				processor.BuildEnv(
					os.Environ(),
					s.EnvPass,
					append(
						s.EnvForce,
						"BOT_SERVER_PORT="+strconv.FormatInt(s.Port, 10),
						"BOT_VERSION="+version,
						"BOT_BUILD_REV="+buildRev,
						"BOT_BUILD_DATE="+buildDate,
					)),
				s.ReplayToUser,
				s.Timeout,
			)
		}
		go sender.Sender(log.WithArea(k+":sender"), s.Token, outgoingQueue)
		go httpif.HttpIf(log.WithArea(k+":http"), s.Port, outgoingQueue)
	}
	<-make(chan interface{})
}
