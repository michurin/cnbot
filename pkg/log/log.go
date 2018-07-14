package log

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Formatter struct{}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	area_name := ""
	switch area := entry.Data["area"].(type) {
	case string:
		area_name = area
	}
	return []byte(
		entry.Time.Format(time.RFC3339) +
			" " + entry.Level.String() + " [" + area_name + "] " + entry.Message + "\n"), nil
}

type Logger struct {
	*logrus.Entry
}

func New() *Logger {
	return &Logger{logrus.NewEntry(&logrus.Logger{
		Out:       os.Stderr,
		Formatter: &Formatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	})}
}

func (l *Logger) WithArea(a string) *Logger {
	return &Logger{l.WithField("area", a)}
}
