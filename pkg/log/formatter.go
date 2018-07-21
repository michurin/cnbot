package log

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Formatter struct {
	isTerminal bool
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	area_name := ""
	switch area := entry.Data["area"].(type) {
	case string:
		area_name = area
	}
	color_level := ""
	color_off := ""
	if f.isTerminal {
		color_off = "\033[0m"
		switch entry.Level {
		case logrus.PanicLevel:
			color_level = "\033[1;33;41m"
		case logrus.FatalLevel:
			color_level = "\033[1;31;43m"
		case logrus.ErrorLevel:
			color_level = "\033[1;31m"
		case logrus.WarnLevel:
			color_level = "\033[1;33m"
		case logrus.InfoLevel:
			color_level = "\033[32m"
		case logrus.DebugLevel:
			color_level = "\033[34m"
		}
	}
	return []byte(
		entry.Time.Format(time.RFC3339) +
			" " + color_level + entry.Level.String() + color_off +
			" [" + area_name + "] " +
			entry.Message + "\n"), nil
}
