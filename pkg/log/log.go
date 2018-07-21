package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Entry
}

func New() *Logger {
	var formatter logrus.Formatter
	var stream *os.File
	var level logrus.Level
	if _, ok := os.LookupEnv("BOT_LOG_STDOUT"); ok {
		stream = os.Stdout
	} else {
		stream = os.Stderr
	}
	if _, ok := os.LookupEnv("BOT_LOG_JSON"); ok {
		formatter = new(logrus.JSONFormatter)
	} else {
		formatter = &Formatter{checkIfTerminal(stream)}
	}
	if _, ok := os.LookupEnv("BOT_LOG_NODEBUG"); ok {
		level = logrus.InfoLevel
	} else {
		level = logrus.DebugLevel
	}
	return &Logger{logrus.NewEntry(&logrus.Logger{
		Out:       stream,
		Formatter: formatter,
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	})}
}

func (l *Logger) WithArea(a string) *Logger {
	return &Logger{l.WithField("area", a)}
}
