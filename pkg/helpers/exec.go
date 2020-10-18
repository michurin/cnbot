package helpers

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"syscall"
	"time"
)

func killGrp(ctx context.Context, pid int, sig syscall.Signal) {
	// We do not consider error as critical because the process could
	// disappear by its own. It is not easy to identify error in this case.
	// For example you can get ESRCH (0x3) that doesn't support by syscall.Errno.Is().
	pgid, err := syscall.Getpgid(pid) // not cmd.SysProcAttr.Pgid
	if err != nil {
		Log(ctx, err)
		return
	}
	err = syscall.Kill(-pgid, sig) // minus
	if err != nil {
		Log(ctx, err)
		return
	}
}

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
	[]byte,
	error,
) {
	Log(ctx, env, command, args)
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
	err := cmd.Start()
	if err != nil {
		Log(ctx, err)
		return nil, err
	}
	// wait command with care about timeouts and ctx
	sync := make(chan error)
	go func() {
		sync <- cmd.Wait()
	}()
	termBound := time.After(termTimeout)
	killBound := time.After(termTimeout + killTimeout)
	waitBound := time.After(termTimeout + killTimeout + waitTimeout)
	for {
		select {
		case err := <-sync: // it has to appear before kill sections to catch stat errors
			if err != nil { // *exec.ExitError if status != 0
				Log(ctx, err, errBuffer.Bytes())
				return nil, err
			}
			if len(errBuffer.Bytes()) != 0 { // just log stderr if any
				Log(ctx, command, args, errBuffer.Bytes())
			}
			if !cmd.ProcessState.Exited() {
				panic("The program is not exited! It's impossible")
			}
			return outBuffer.Bytes(), nil
		case <-ctx.Done(): // urgent exit, we doesn't even wait for process finalization
			Log(ctx, "Exec terminated by context")
			killGrp(ctx, cmd.Process.Pid, syscall.SIGKILL)
			return nil, nil
		case <-termBound:
			killGrp(ctx, cmd.Process.Pid, syscall.SIGTERM)
		case <-killBound:
			killGrp(ctx, cmd.Process.Pid, syscall.SIGKILL)
		case <-waitBound:
			// Very bad case, we gave up, we leave goroutine with cmd.Wait running...
			// I hope it will never happen.
			err := errors.New("can't wait anymore")
			Log(ctx, err)
			return nil, err
		}
	}
}
