package bot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/michurin/minlog"
)

var /* const */ (
	startedAtTime = time.Now()
	startedAt     = startedAtTime.Format(time.RFC3339)
	startedAtUnix = startedAtTime.Unix()
	pid           = os.Getpid()
	goVer         = runtime.Version()
)

type AliveHandler struct{}

func memMap(m *runtime.MemStats) map[string]interface{} {
	v := reflect.ValueOf(*m)
	t := reflect.TypeOf(*m)
	e := map[string]interface{}{}
	for i := 0; i < v.NumField(); i++ {
		switch t.Field(i).Type.Kind() { // nolint:exhaustive
		case reflect.Bool, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			e[t.Field(i).Name] = v.Field(i).Interface()
		default:
		}
	}
	return e
}

func (h *AliveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	b, err := json.Marshal(map[string]interface{}{
		"version": map[string]interface{}{
			"version": version,
			"build":   Build,
			"go":      goVer,
			"goos":    runtime.GOOS,
			"goarch":  runtime.GOARCH,
		},
		"started_at":      startedAt,
		"started_at_unix": startedAtUnix,
		"pid":             pid,
		"num_goroutine":   runtime.NumGoroutine(),
		"num_cpu":         runtime.NumCPU(),
		"mem_status":      memMap(m),
	})
	if err != nil {
		minlog.Log(r.Context(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	n, err := w.Write(b)
	if err != nil {
		minlog.Log(r.Context(), err, fmt.Sprintf("(%d/%d written)", n, len(b)))
		return
	}
}
