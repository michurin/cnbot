package execute

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"

	"github.com/michurin/cnbot/pkg/interfaces"
	"github.com/pkg/errors"
)

type Executor struct {
	Timeout    time.Duration
	Logger     interfaces.Logger
	Command    string
	Cwd        string
	KillSignal syscall.Signal
}

func New(logger interfaces.Logger) *Executor {
	return &Executor{
		Timeout:    2 * time.Second, // TODO
		Logger:     logger,
		Command:    "./test_user_script.sh", // TODO it seems it have to come to Run from Task
		Cwd:        ".",                     // TODO ...and it too
		KillSignal: syscall.SIGKILL,
	}
}

func (e *Executor) Run(ctx context.Context, env []string, args []string) ([]byte, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	// CommandContext still do not use GPID, so I need to control timeout manually
	// to avoid problems with process that spawn children. Like sh script that does sleep(1)
	cmd := exec.Command(e.Command, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // setpgid(2) between fork(2) and execve(2)
	cmd.Env = env
	cmd.Dir = e.Cwd
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	e.Logger.Log(fmt.Sprintf("Run %+v", cmd))
	err := cmd.Start()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ctx, cancel := context.WithTimeout(ctx, e.Timeout)
	defer cancel()
	var processAlreadyDone bool
	go func() {
		<-ctx.Done()
		if processAlreadyDone {
			return
		}
		e.Logger.Log("Kill process due to timeout")
		// cmd.Process is not nil here because we are started
		// -PID is the same as -PGID
		// https://medium.com/@felixge/killing-a-child-process-and-all-of-its-children-in-go-54079af94773
		err := syscall.Kill(-cmd.Process.Pid, e.KillSignal)
		if err != nil {
			e.Logger.Log(errors.WithStack(err))
		}
	}()
	err = cmd.Wait()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	processAlreadyDone = true
	// TODO check error code
	// TODO check stderr
	return stdout.Bytes(), nil
}
