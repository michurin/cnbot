package processor

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/michurin/cnbot/pkg/log"
	"github.com/michurin/cnbot/pkg/prepareoutgoing"
	"github.com/michurin/cnbot/pkg/receiver"
	"github.com/michurin/cnbot/pkg/sender"
)

func BuildEnv(env []string, envPass []string, envForce []string) []string {
	res := []string{}
	allowed := map[string]bool{}
	for _, k := range envPass {
		allowed[k] = true
	}
	for _, kv := range env {
		if allowed[strings.SplitN(kv, "=", 2)[0]] {
			res = append(res, kv)
		}
	}
	return append(append(res, "BOT_PID="+strconv.Itoa(os.Getpid())), envForce...)
}

func Processor(
	log *log.Logger,
	inQueue <-chan receiver.TUpdateMessage,
	outQueue chan<- sender.OutgoingData,
	whitelist []int64,
	command string,
	cwd string,
	env []string,
	timeout int64,
) {
	for message := range inQueue {
		if intInSlice(message.From.Id, whitelist) {
			outData := execute(
				log,
				command,
				cwd,
				append(
					env,
					"BOT_USER_NAME="+message.From.Username,
					"BOT_USER_ID="+strconv.FormatInt(message.From.Id, 10),
					"BOT_CHAT_ID="+strconv.FormatInt(message.Chat.Id, 10),
				),
				timeout,
				strings.Fields(message.Text), // TODO: make it configurable?
			)
			q := prepareoutgoing.PrepareOutgoing(log, outData, message.From.Id, nil)
			if q.MessageType != "" {
				outQueue <- q
			}
		} else {
			log.Infof("WARNING: from_id=%d is not allowed. Add to whitelist", message.From.Id)
			outQueue <- prepareoutgoing.PrepareOutgoing(
				log,
				[]byte(fmt.Sprintf("Sorry. Your ID (%d) is not allowd.", message.From.Id)),
				message.From.Id,
				nil,
			)
			continue
		}
	}
}
