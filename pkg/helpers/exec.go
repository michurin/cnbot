package helpers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

// Note: don't use ctx for timeouts
func Exec(
	ctx context.Context,
	termTimeout time.Duration,
	killTimeout time.Duration,
	waitTimeout time.Duration,
	command string,
	args []string,
	env []string,
	pwd string,
) (
	stdout []byte,
	stderr []byte,
	err error,
) {
	// setup cmd
	cmd := exec.Command(command, args...) // we don't use ctx here because it kills only process, not group
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Dir = pwd
	cmd.Env = env
	var outBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	var errBuffer bytes.Buffer
	cmd.Stderr = &errBuffer
	// start command synchronously; we hope it doesn't take a long time
	err = cmd.Start()
	if err != nil {
		Log(ctx, err)
		return nil, nil, err
	}
	// wait command with care about timeouts and ctx
	sync := make(chan error)
	go func() {
		sync <- cmd.Wait()
	}()
	pgid, err := syscall.Getpgid(cmd.Process.Pid) // not cmd.SysProcAttr.Pgid
	if err != nil {
		Log(ctx, err)
		return nil, nil, err
	}
	pgid = -pgid // minus!
	termBound := time.After(termTimeout)
	killBound := time.After(termTimeout + killTimeout)
	waitBound := time.After(termTimeout + killTimeout + waitTimeout)
	for {
		select {
		case err := <-sync: // it has to appear before kill sections to catch stat errors
			if err != nil {
				Log(ctx, err)
				return nil, nil, err
			}
			if !cmd.ProcessState.Exited() {
				panic("The program is not exited! It's impossible")
			}
			exitCode := cmd.ProcessState.ExitCode()
			if exitCode != 0 {
				return nil, nil, fmt.Errorf("Exit code %d: %s", exitCode, string(stderr))
			}
			return outBuffer.Bytes(), errBuffer.Bytes(), nil
		case <-ctx.Done(): // we doesn't even wait for process finalization
			err = syscall.Kill(pgid, syscall.SIGKILL)
			if err != nil {
				Log(ctx, err)
				return nil, nil, err
			}
			return nil, nil, errors.New("aborted by context")
		case <-termBound:
			err = syscall.Kill(pgid, syscall.SIGTERM)
			if err != nil {
				Log(ctx, err)
				return nil, nil, err
			}
		case <-killBound:
			err = syscall.Kill(pgid, syscall.SIGKILL)
			if err != nil {
				Log(ctx, err)
				return nil, nil, err
			}
		case <-waitBound:
			err := errors.New("can't wait anymore")
			Log(ctx, err)
			return nil, nil, err
		}
	}
}
