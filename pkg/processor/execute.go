package processor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/michurin/cnbot/pkg/log"
)

func execute(
	ctx context.Context,
	log *log.Logger,
	command string,
	cwd string,
	env []string,
	args []string,
) []byte {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var outData []byte
	cmd := exec.Command(command, args...)                 // CommandContext still do not use GPID, so we control timeout manually
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // setpgid(2) between fork(2) and execve(2)
	cmd.Env = env
	cmd.Dir = cwd
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	log.Infof("Run %+v", cmd)
	err := cmd.Start()
	if err == nil {
		go func() {
			<-ctx.Done()
			// cmd.Process is not nil here cose we are started
			// -PID is the same as -PGID
			// https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}()
		err = cmd.Wait()
	}
	if err == nil {
		log.Info("Done.")
		outData = stdout.Bytes()
	} else {
		log.Warnf("Done with error %v", err)
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
