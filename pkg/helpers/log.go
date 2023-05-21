package helpers

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"os"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/michurin/jsonpainter"
	"github.com/michurin/minlog"
)

func SetupLogging() {
	_, fileName, _, ok := runtime.Caller(0)
	if !ok {
		panic("can not get build prefix")
	}
	buildPrefixLen := len(fileName) - len("helpers/log.go")
	callerOpt := minlog.WithCallerCutter(func(p string) string {
		return p[buildPrefixLen:]
	})
	fi, err := os.Stdout.Stat()
	if err != nil {
		minlog.Log(context.Background(), "Can not check stat(stdout)", err)
		return
	}
	if (fi.Mode() & os.ModeCharDevice) != 0 {
		jcOpts := []jsonpainter.Option{
			jsonpainter.ClrCtl(jsonpainter.Blue),
			jsonpainter.ClrKey(jsonpainter.Brown),
		}
		minlog.SetDefaultLogger(minlog.New(
			callerOpt,
			minlog.WithLineFormatter(func(tm, level, label, caller, msg string) string {
				lc := "1"
				if level == minlog.DefaultInfoLabel {
					lc = "2"
				}
				return strings.Join([]string{
					tm,
					"\033[3" + lc + "m" + level + "\033[0m",
					"\033[1m" + label + "\033[0m",
					"\033[34m" + caller + "\033[0m",
					jsonpainter.String(msg, jcOpts...),
				}, " ")
			}),
		))
	} else {
		minlog.SetDefaultLogger(minlog.New(
			callerOpt,
			minlog.WithLineFormatter(func(tm, level, label, caller, msg string) string {
				return strings.Join([]string{
					tm,
					"[" + level + "]",
					label,
					"[" + caller + "]",
					msg,
				}, " ")
			}),
		))
	}
}

var labelCounter uint32

func AutoLabel() string {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, atomic.AddUint32(&labelCounter, 1))
	return hex.EncodeToString(b[2:4])
}
