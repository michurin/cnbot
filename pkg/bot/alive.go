package bot

import (
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"time"

	hps "github.com/michurin/cnbot/pkg/helpers"
)

var /* const */ startedAt = time.Now().Format(time.RFC3339)
var /* const */ pid = os.Getpid()
var /* const */ goVer = runtime.Version()

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
		"started_at":    startedAt,
		"pid":           pid,
		"num_goroutine": runtime.NumGoroutine(),
		"num_cpu":       runtime.NumCPU(),
		"mem_status":    memMap(m),
	})
	if err != nil {
		hps.Log(r.Context(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	n, err := w.Write(b)
	if err != nil {
		hps.Log(r.Context(), err)
		return
	}
	if n != len(b) {
		hps.Log(r.Context(), "Not all data has been written")
		return
	}
}
