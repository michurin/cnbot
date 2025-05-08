package apiserver

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type APIAct struct {
	IsJSON   bool // TODO use content type?
	Stream   bool // TODO couple with IsJSON
	Request  string
	Response []byte
}

func APIServer(t *testing.T, cancel context.CancelFunc, api map[string][]APIAct) (string, func()) {
	// TODO
	// HUGE MISFEATURE: we do not check all calls are really happened! We are checking only happened calls though.
	t.Helper()
	testDone := make(chan struct{})
	steps := map[string]int{} // it looks ugly, however we can use it without locks
	tg := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// DO NOT user require.* in this handler.
		// require.* is based on t.FailNow() it won't work in goroutines
		// assert.* is founded on t.Fatal()
		// assert.Equal(t, http.MethodPost, r.Method) // TODO! Assert method
		bodyBytes, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		body := string(bodyBytes)

		url := r.URL.String()
		t.Logf("Mock server: url=%q", url)
		n := steps[url]
		ax, ok := api[url]
		assert.True(t, ok, "URL not found: "+url)
		a := ax[n]
		steps[url] = n + 1
		if a.Stream { // TODO couple Stream and IsJSON
			_, err = w.Write(a.Response)
			assert.NoError(t, err)
			return
		}
		if a.IsJSON {
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.JSONEq(t, a.Request, body)
		} else {
			ctype := r.Header.Get("Content-Type")
			assert.Contains(t, ctype, "multipart/form-data")
			idx := strings.Index(ctype, "boundary=")
			assert.Greater(t, idx, -1, "ctype="+ctype)
			universal := strings.ReplaceAll(body, ctype[idx+9:], "BOUND")
			assert.Equal(t, a.Request, universal)
		}
		if a.Response == nil {
			cancel()
			<-testDone
		}
		_, err = w.Write(a.Response)
		assert.NoError(t, err)
	}))
	return tg.URL, func() {
		close(testDone)
		tg.Close()
	}
}
