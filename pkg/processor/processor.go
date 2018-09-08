package processor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/michurin/cnbot/pkg/log"
	"github.com/michurin/cnbot/pkg/prepareoutgoing"
	"github.com/michurin/cnbot/pkg/receiver"
	"github.com/michurin/cnbot/pkg/sender"
)

var (
	errorNoFrom = errors.New("Can not find <from> info")
	errorNoChat = errors.New("Can not find <chat> info")
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

func messageToFrom(p receiver.TUpdateResult) (receiver.TUpdateFrom, error) {
	if p.Message != nil {
		if p.Message.From != nil {
			return *p.Message.From, nil
		}
	} else if p.CallbackQuery != nil {
		return p.CallbackQuery.From, nil
	}
	return receiver.TUpdateFrom{}, errorNoFrom
}

func getChat(p receiver.TUpdateResult) (receiver.TUpdateChat, error) {
	if p.Message != nil {
		return p.Message.Chat, nil
	}
	// chat_id not present in callback_query messages
	return receiver.TUpdateChat{}, errorNoChat
}

func Processor(
	log *log.Logger,
	inQueue <-chan receiver.TUpdateResult,
	outQueue chan<- sender.OutgoingData,
	whitelist []int64,
	command string,
	cwd string,
	env []string,
	replayToUser bool,
	ma MessageToArgs,
	timeout int64,
) {
	for part := range inQueue {
		from, err := messageToFrom(part)
		if err != nil {
			log.Warnf("%s: %+v", err, part)
			continue
		}
		chat, err := getChat(part)
		chatId := chat.Id
		if err != nil {
			log.Infof("Force chat_id=user_id: %s", err)
			chatId = from.Id
		}
		targetId := chatId
		if replayToUser {
			targetId = from.Id
		}
		if intInSlice(targetId, whitelist) {
			isCallBack, args, err := ma.MessageToArgs(part)
			if err != nil {
				log.Warnf("%s: %+v", err, part)
				continue
			}
			if isCallBack {
				q, err := prepareoutgoing.CallbackAnswerOutgoing(part.CallbackQuery.Id)
				if err != nil {
					log.Error("Can not prepare CallbackAnswerOutgoing: %s", err.Error())
					continue
				}
				outQueue <- q
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
			outData := execute(
				ctx,
				log,
				command,
				cwd,
				append(
					env,
					"BOT_USER_NAME="+from.Username,
					"BOT_USER_ID="+strconv.FormatInt(from.Id, 10),
					"BOT_CHAT_ID="+strconv.FormatInt(chatId, 10),
					"BOT_TARGET_ID="+strconv.FormatInt(targetId, 10),
				),
				args,
			)
			cancel()
			q, err := prepareoutgoing.PrepareOutgoing(
				log,
				outData,
				targetId,
				nil,
				nil,
			)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			if q.MessageType != "" {
				outQueue <- q
			}
		} else {
			log.Warnf(
				"WARNING: from_id=%d is not allowed. Add to whitelist",
				targetId,
			)
			q, err := prepareoutgoing.PrepareOutgoing(
				log,
				[]byte(fmt.Sprintf(
					"Sorry. Your effective ID %d not allowed.",
					targetId,
				)),
				targetId,
				nil,
				nil,
			)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			outQueue <- q
		}
	}
}
