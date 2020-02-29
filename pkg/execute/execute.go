package execute

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/michurin/cnbot/pkg/interfaces"
	"github.com/pkg/errors"
)

type ScriptInfo struct {
	Name    string
	Timeout time.Duration
	Env     []string
	Args    []string
}

type Executor struct {
	Logger     interfaces.Logger // TODO make it private?
	KillSignal syscall.Signal
	Env        []string
}

func New(logger interfaces.Logger, commonEnv []string) *Executor {
	return &Executor{
		Logger:     logger,
		KillSignal: syscall.SIGKILL,
		Env:        commonEnv,
	}
}

func (e *Executor) Run(ctx context.Context, script ScriptInfo) ([]byte, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command, err := filepath.Abs(script.Name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// CommandContext still do not use GPID, so I need to control timeout manually
	// to avoid problems with process that spawn children. Like sh script that does sleep(1)
	cmd := exec.Command(command, script.Args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // setpgid(2) between fork(2) and execve(2)
	cmd.Env = append(e.Env, script.Env...)
	cmd.Dir = path.Dir(command) // TODO configurable?
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	e.Logger.Log(fmt.Sprintf("Run %v %+v", cmd.Env, cmd))
	err = cmd.Start()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ctx, cancel := context.WithTimeout(ctx, script.Timeout)
	defer cancel()
	go func() {
		<-ctx.Done()
		if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
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
	exitCode := cmd.ProcessState.ExitCode()
	errOutput := stderr.String()
	if exitCode != 0 || len(errOutput) > 0 {
		return nil, errors.New(fmt.Sprintf("exitCode=%d stderr=\"%s\"", exitCode, errOutput))
	}
	return stdout.Bytes(), nil
}
