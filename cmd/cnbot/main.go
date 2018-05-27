package main

// bot name @M_78c9409716d3f0bdfd6d_bot

import (
//	"github.com/michurin/cnbot/pkg/cnbot"
	"github.com/michurin/cnbot/pkg/config"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"bytes"
	"io/ioutil"
)

type handler struct {
	cfg config.BotConfiguration
}

func NewHandler(c config.BotConfiguration) handler {
	return handler{cfg: c}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!\nConfig: %#v\n", r.URL.Path[1:], h.cfg)
}

// ----------------------------------

type Client struct {
	httpClient *http.Client
	apiURL string
}

func (c Client) PostJSON(data []byte) ([]byte, error) {
	resp, err := c.httpClient.Post(c.apiURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

type ClientFactory struct {
	Proxy *url.URL
	Token string
	Timeout time.Duration
}

func (f ClientFactory) NewClient(apiMethod string) Client {
	return Client{
		httpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: func(_ *http.Request) (*url.URL, error) {  // TODO setup proxy only if proxy configured
					return f.Proxy, nil
				},
			},
			Timeout: f.Timeout,
		},
		apiURL: "https://api.telegram.org/bot" + f.Token + "/" + apiMethod,
	}
}

// ----------------------------------

// TODO NEXT STEPS
// 1. Change /getMe to /getUpdates (polling loop in background)
// 2. Spawn processes
// 3. /sendMessage
// 4. Write custom wrapper for http.Client
//    - move proxy, timeout, token[?]
//    - methods .PostJson, .PostMultipart
func StartPollingLoop(client Client, pollingInterfal int, q chan<- int) {
	buf, err := client.PostJSON([]byte("{}"))
	if err != nil {
		fmt.Println("ERROR", err)
		return
	}
	fmt.Println(string(buf))
}

func StartProcessor(q <-chan int, r chan<- bool) {}
func StartServer(r chan<- bool) {}
func StartResponder(r <-chan bool) {}

func StartBot(cfg config.BotConfiguration) {
	clientFactory := ClientFactory{
		Proxy: cfg.Proxy,
		Token: cfg.Token,
		Timeout: cfg.Timeout,
	}
	client_message_queue := make(chan int, 1000)
	server_message_queue := make(chan bool, 1000)
	StartPollingLoop(clientFactory.NewClient("getMe"), cfg.PollingInterval, client_message_queue)
	StartProcessor(client_message_queue, server_message_queue)  // предзапустить
	StartServer(server_message_queue)
	StartResponder(server_message_queue)
}

func main() {
	botconfigs := config.GetConfiguration()
	fmt.Printf("%#v\n", botconfigs)
	StartBot(botconfigs[0])
	return

	s := &http.Server{
		Addr:              ":8080",
		Handler:           NewHandler(botconfigs[0]),
		ReadHeaderTimeout:  1 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		MaxHeaderBytes:    1 << 12,
	}
	fmt.Println("Start...")
	err := s.ListenAndServe()
	fmt.Println(err)
//	cnbot.PollingLoop(proxyServer)
}
