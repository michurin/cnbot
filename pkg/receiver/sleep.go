package receiver

import (
	"time"

	"github.com/michurin/cnbot/pkg/log"
)

func errorSleep(log *log.Logger, s int) {
	log.Warnf("Error sleep %d seconds", s)
	time.Sleep(time.Duration(s) * time.Second)
}
