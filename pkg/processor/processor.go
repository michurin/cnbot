package processor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/michurin/cnbot/pkg/log"
	"github.com/michurin/cnbot/pkg/prepareoutgoing"
	"github.com/michurin/cnbot/pkg/receiver"
	"github.com/michurin/cnbot/pkg/sender"
)

func execute(
	log *log.Logger,
	command string,
	cwd string,
	env []string,
	timeout int64,
	message string,
	fromId int64,
) []byte {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var outData []byte
	cmd := exec.Command(command, message)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // setpgid(2) between fork(2) and execve(2)
	cmd.Env = append(env, "BOT_CHAT_ID="+strconv.FormatInt(fromId, 10))
	cmd.Dir = cwd
	stdout.Reset()
	stderr.Reset()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	afterFuncTimer := time.AfterFunc(time.Second*time.Duration(timeout), func() { // TODO defer clean timer
		// https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773
		if cmd.Process != nil { // nil if script not started (not exists)
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) // -PID is the same as -PGID
		}
	})
	defer afterFuncTimer.Stop() // cosmetics, all timers share the same gorutine
	log.Info("Run...")
	err := cmd.Run()
	log.Info("Done.")
	if err == nil {
		outData = stdout.Bytes()
	} else {
		outData = []byte("Subprocess error") // Default message
		switch v := err.(type) {
		case *exec.ExitError:
			pid := v.ProcessState.Pid()
			status, ok := v.Sys().(syscall.WaitStatus)
			if ok {
				code := status.ExitStatus()
				stopped := status.Signaled()
				log.Infof("%T\nPID=%d\nCode=%d\nSignaled=%v", v, pid, code, stopped)
				outData = []byte(fmt.Sprintf("Subprocess error: Code=%d, Signaled=%t, stderr=%s", code, stopped, stderr.String()))
			} else {
				log.Infof("Can not parse status %#v", err)
				outData = []byte("Subprocess error: can not parse status")
			}
		case *os.PathError:
			log.Infof("Command not found: %s", v.Path)
			outData = []byte(fmt.Sprintf("Subprocess error: Command %s not found", v.Path))
		default:
			log.Info("UNKNOWN ERROR")
			log.Infof("%T\n%#v", err, err)
			log.Info(fmt.Sprint(err) + ": " + stderr.String())
		}
	}
	return outData
}

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
	for {
		message := <-inQueue
		if intInSlice(message.From.Id, whitelist) {
			outData := execute(log, command, cwd, env, timeout, message.Text, message.From.Id)
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
