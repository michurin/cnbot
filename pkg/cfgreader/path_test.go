package cfgreader

import (
	"fmt"
	"testing"
)

func TestPath(t *testing.T) {
	for i, c := range []struct {
		cfg      string
		script   string
		expected string
	}{
		{"ok.ini", "", ""},
		{"ok.ini", "/abs/path/ok.sh", "/abs/path/ok.sh"},
		{"ok.ini", "rel/ok.sh", "rel/ok.sh"},
		{"./ok.ini", "rel/ok.sh", "rel/ok.sh"},
		{"../ok.ini", "rel/ok.sh", "../rel/ok.sh"},
		{"../ok.ini", "../rel/../ok.sh", "../../ok.sh"},
	} {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			if pathToScript(c.cfg, c.script) != c.expected {
				t.Fail()
			}
		})
	}
}
