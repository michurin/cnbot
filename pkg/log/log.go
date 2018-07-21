package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Entry
}

func New() *Logger {
	stream := os.Stderr
	return &Logger{logrus.NewEntry(&logrus.Logger{
		Out:       stream,
		Formatter: &Formatter{checkIfTerminal(stream)},
		Level:     logrus.DebugLevel,
	})}
}

func (l *Logger) WithArea(a string) *Logger {
	return &Logger{l.WithField("area", a)}
}
