package helpers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
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

var buildPrefixLen int
var labelInfo = "[info]"
var labelError = "[error]"
var labelPrefix = "["
var labelPostfix = "]"
var callerPrefix = "["
var callerPostfix = "]"

func fmtMessage(messages ...interface{}) (label, msg string) {
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
		case time.Duration:
			msg = m.String()
		case []string:
			msg = "[" + strings.Join(m, ", ") + "]"
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
	msg = strings.Join(msgs, " ")
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
	level, msg := fmtMessage(message...)
	fmt.Printf("%s %s %s%s%s %s%s:%d%s %s\n", tm, level, labelPrefix, label, labelPostfix, callerPrefix, file[buildPrefixLen:], line, callerPostfix, msg)
}

func Label(ctx context.Context, labels ...string) context.Context {
	label := strings.Join(labels, ":")
	prevLabel, ok := ctx.Value(labelKey).(string)
	if ok {
		label = prevLabel + ":" + label
	}
	return context.WithValue(ctx, labelKey, label)
}

func RandLabel() string {
	b := make([]byte, 3)
	n, err := rand.Reader.Read(b)
	if err != nil {
		Log(context.Background(), err)
	}
	if n != len(b) {
		Log(context.Background(), errors.New("not enough data"))
	}
	return hex.EncodeToString(b)
}
