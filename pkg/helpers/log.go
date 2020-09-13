package helpers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func init() {
	_, fileName, _, ok := runtime.Caller(0)
	if !ok {
		panic("can not get build prefix")
	}
	buildPrefixLen = len(fileName) - len("pkg/helpers/log.go")
	fi, err := os.Stdout.Stat()
	if err != nil {
		Log(context.Background(), "Can not check stat(stdout)", err)
		return
	}
	if (fi.Mode() & os.ModeCharDevice) != 0 {
		labelInfo = "\033[32minfo\033[0m"
		labelError = "\033[31merror\033[0m"
	}
}

var buildPrefixLen int
var labelInfo = "info"
var labelError = "error"

func fmtMessage(messages ...interface{}) (label, msg string) {
	label = labelInfo
	msgs := []string(nil)
	for _, message := range messages {
		switch m := message.(type) {
		case string:
			msg = m
		case []byte:
			msg = string(m)
		case int:
			msg = fmt.Sprintf("%d", m)
		case nil:
			msg = "<nil>"
		case *exec.ExitError:
			label = labelError
			msg = fmt.Sprintf("Exit error [code=%d]: %s: %s", m.ExitCode(), m.Error(), string(m.Stderr))
		case error:
			label = labelError
			msg = m.Error()
		default:
			label = labelError
			msg = fmt.Sprintf("Unknown type [%[1]T]: %+[1]v", m)
		}
		msgs = append(msgs, msg)
	}
	msg = strings.Join(msgs, " ")
	return
}

// Crucial ideas
// - to be extremely simple to be used and mocked
// - use context to track requests id and similar labels
// - derive log level from type of argument
// - smart formatting
// - print caller
func Log(ctx context.Context, message ...interface{}) {
	tm := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "[nofile]"
	}
	level, msg := fmtMessage(message...)
	fmt.Printf("%s [%s] [%s:%d] %s\n", tm, level, file[buildPrefixLen:], line, msg)
}
