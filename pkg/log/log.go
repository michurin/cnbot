// Naive implementation of logger. You are free to implement this interface and use your
// favorite logging tool.

package log

import (
	"fmt"
	"log"
	"os"
)

type Logger struct{ l *log.Logger }

func New() *Logger {
	return &Logger{l: log.New(os.Stdout, "", log.Ldate|log.Ltime)}
}

func (l Logger) Log(msg interface{}) {
	l.l.Print(value(msg))
}

func value(msg interface{}) string {
	switch v := msg.(type) {
	case nil:
		return "<nil>"
	case string:
		return v
	case []byte:
		return string(v)
	case interface {
		error
		fmt.Formatter
	}:
		return fmt.Sprintf("%+v", v)
	case error:
		return v.Error()
	default:
		return fmt.Sprintf("%T: %v", v, v)
	}
}
