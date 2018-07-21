package log

import (
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func checkIfTerminal(w io.Writer) bool {
	v, ok := w.(*os.File)
	if ok {
		return terminal.IsTerminal(int(v.Fd()))
	}
	return false
}
