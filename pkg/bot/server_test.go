package bot

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecoder(t *testing.T) {
	t.Parallel()

	t.Run("invalid_method", func(t *testing.T) {
		r := newReq(t, http.MethodGet, "", "", "")
		u, b, c, err := decodeRequest(r)
		assert.Error(t, err)
		assert.Empty(t, b)
		assert.Empty(t, u)
		assert.Empty(t, c)
	})

	for _, url := range []string{"", "http://1.1.1.1/", "http://1.1.1.1/x"} {
		url := url
		t.Run("raw_mode_invalid_url_"+url, func(t *testing.T) {
			r := newReq(t, http.MethodPost, "", url, "")
			u, b, c, err := decodeRequest(r)
			assert.Error(t, err)
			assert.Empty(t, b)
			assert.Empty(t, u)
			assert.Empty(t, c)
		})
	}

	t.Run("raw_ok", func(t *testing.T) {
		r := newReq(t, http.MethodPost, "", "http://1.1.1.1/12", "text")
		u, b, c, err := decodeRequest(r)
		assert.Nil(t, err)
		assert.Equal(t, []byte("text"), b)
		assert.Equal(t, int64(12), u)
		assert.Empty(t, c)
	})

	cntType := "multipart/form-data; boundary=bbb"

	t.Run("multipart_can_not_parse", func(t *testing.T) {
		r := newReq(t, http.MethodPost, cntType, "http://1.1.1.1/", "")
		u, b, c, err := decodeRequest(r)
		assert.Error(t, err)
		assert.Empty(t, b)
		assert.Empty(t, u)
		assert.Empty(t, c)
	})

	body := "--bbb\nContent-Disposition: form-data; name=\"to\"\n\n12\n" +
		"--bbb\nContent-Disposition: form-data; name=\"msg\"\n\ntext\n" +
		"--bbb--\n"

	for _, url := range []string{"http://1.1.1.1/", "http://1.1.1.1/14", "http://1.1.1.1/nan"} {
		url := url
		t.Run("multipart_can_not_parse_"+url, func(t *testing.T) {
			r := newReq(t, http.MethodPost, cntType, url, body)
			u, b, c, err := decodeRequest(r)
			assert.Nil(t, err)
			assert.Equal(t, []byte("text"), b)
			assert.Equal(t, int64(12), u)
			assert.Empty(t, c)
		})
	}

	body = "--bbb\nContent-Disposition: form-data; name=\"to\"\n\n12\n" +
		"--bbb\nContent-Disposition: form-data; name=\"msg\"\n\ntext\n" +
		"--bbb\nContent-Disposition: form-data; name=\"cap\"\n\ncaption\n" +
		"--bbb--\n"

	t.Run("multipart_with_caption", func(t *testing.T) {
		r := newReq(t, http.MethodPost, cntType, "http://1.1.1.1/", body)
		u, b, c, err := decodeRequest(r)
		assert.Nil(t, err)
		assert.Equal(t, []byte("text"), b)
		assert.Equal(t, int64(12), u)
		assert.Equal(t, "caption", c)
	})

	body = "--bbb\nContent-Disposition: form-data; name=\"msg\"\n\ntext\n--bbb--\n"

	t.Run("multipart_without_to", func(t *testing.T) {
		r := newReq(t, http.MethodPost, cntType, "http://1.1.1.1/12", body)
		u, b, c, err := decodeRequest(r)
		assert.Nil(t, err)
		assert.Equal(t, []byte("text"), b)
		assert.Equal(t, int64(12), u)
		assert.Empty(t, c)
	})

	t.Run("multipart_without_to_without_fallback", func(t *testing.T) {
		r := newReq(t, http.MethodPost, cntType, "http://1.1.1.1/", body)
		u, b, c, err := decodeRequest(r)
		assert.Error(t, err)
		assert.Empty(t, b)
		assert.Empty(t, u)
		assert.Empty(t, c)
	})

	body = "--bbb\nContent-Disposition: form-data; name=\"to\"\n\n12\n--bbb--\n"

	t.Run("multipart_without_msg", func(t *testing.T) {
		r := newReq(t, http.MethodPost, cntType, "http://1.1.1.1/", body)
		u, b, c, err := decodeRequest(r)
		assert.Error(t, err)
		assert.Empty(t, b)
		assert.Empty(t, u)
		assert.Empty(t, c)
	})

	t.Run("multipart_without_msg", func(t *testing.T) {
		r := newReq(t, http.MethodPost, cntType, "http://1.1.1.1/", body)
		u, b, c, err := decodeRequest(r)
		assert.Error(t, err)
		assert.Empty(t, b)
		assert.Empty(t, u)
		assert.Empty(t, c)
	})
}

func newReq(t *testing.T, method string, cntType string, url string, body string) *http.Request {
	r, err := http.NewRequestWithContext(context.Background(), method, url, strings.NewReader(body))
	require.Nil(t, err)
	if cntType != "" {
		r.Header.Set("Content-Type", cntType)
	}
	return r
}
