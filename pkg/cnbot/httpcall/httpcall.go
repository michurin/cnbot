package httpcall

import (
	"net/http"
	"golang.org/x/net/proxy"
	"github.com/michurin/cnbot/pkg/perror"
	"bytes"
	"io/ioutil"
	"log"
)

type HttpRequestDoer interface { // TODO не экспортится; можно мокать, проксить и вообще
	Do(*http.Request) (*http.Response, error)
}

type TAPIClient struct {  // TODO не экспортится
	httpClient *http.Client
}

func NewTAPIClient() (*TAPIClient, error) {
	dialer, err := proxy.SOCKS5("tcp", "localhost:8059", nil, proxy.Direct)
	if err != nil {
		return nil, perror.NewErrorString("socks5 creation error: %s", err)
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.Dial = dialer.Dial // TODO use modern DealContext?
	return &TAPIClient{httpClient}, nil
}

func (c TAPIClient) Do(r *http.Request) (*http.Response, error) {
	return c.httpClient.Do(r)
}

// TODO тесты, с подменой интерфейса
func HTTPCall(c HttpRequestDoer, method string, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	if err != nil {
		return nil, perror.NewErrorString("can't create request:", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, perror.NewErrorString("can't GET page:", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, perror.NewErrorString("error reading body:", err)
	}
	log.Printf("%s -> %s", body, b)
	return b, nil
}
