package helpers

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
	"unicode/utf8"
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
		labelPrefix = "\033[1m"
		labelPostfix = "\033[0m"
		callerPrefix = "\033[34m"
		callerPostfix = "\033[0m"
	}
}

type logContextKey string

const labelKey = logContextKey("label")

var (
	buildPrefixLen int
	labelInfo      = "[info]"
	labelError     = "[error]"
	labelPrefix    = "["
	labelPostfix   = "]"
	callerPrefix   = "["
	callerPostfix  = "]"
)

func fmtMessage(sep string, messages ...interface{}) (label, msg string) {
	label = labelInfo
	msgs := []string(nil)
	for _, message := range messages {
		switch m := message.(type) {
		case string:
			msg = m
		case []byte:
			if utf8.Valid(m) {
				msg = string(m)
			} else {
				if len(m) > 200 {
					m = m[:200]
				}
				msg = fmt.Sprintf("%q", m)
			}
		case int:
			msg = fmt.Sprintf("%d", m)
		case int64:
			msg = Itoa(m)
		case time.Duration:
			msg = m.String()
		case []string:
			msg = strings.Join(m, " ")
		case nil:
			msg = "<nil>"
		case error:
			label = labelError
			msg = m.Error()
		default:
			label = labelError
			msg = fmt.Sprintf("Unknown type [%[1]T]: %+[1]v", m)
		}
		msgs = append(msgs, msg)
	}
	msg = strings.Join(msgs, sep)
	return
}

// Crucial ideas
// - to be extremely simple to be used and mocked
// - use context to track requests id and similar labels
// - derive log level from type of argument
// - smart formatting
// - print caller
// TODO
// - rewrite it to struct and use Logger from log
// - inject it everywhere
func Log(ctx context.Context, message ...interface{}) {
	tm := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "[nofile]"
	}
	label, ok := ctx.Value(labelKey).(string)
	if !ok {
		label = "root"
	}
	err := ctx.Err()
	if err != nil {
		message = append(message, "") // 3 lines: prepend
		copy(message[1:], message)
		message[0] = "[ctxErr=" + err.Error() + "]"
	}
	level, msg := fmtMessage(" ", message...)
	fmt.Printf("%s %s %s%s%s %s%s:%d%s %s\n", tm, level, labelPrefix, label, labelPostfix, callerPrefix, file[buildPrefixLen:], line, callerPostfix, msg)
}

func Label(ctx context.Context, labels ...interface{}) context.Context {
	_, label := fmtMessage(":", labels...)
	prevLabel, ok := ctx.Value(labelKey).(string)
	if ok {
		label = prevLabel + ":" + label
	}
	return context.WithValue(ctx, labelKey, label)
}

var labelCounter uint32

func RandLabel() string {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, atomic.AddUint32(&labelCounter, 1))
	return hex.EncodeToString(b[2:4])
}
