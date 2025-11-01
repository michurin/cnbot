package tests_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/michurin/cnbot/pkg/ctxlog"
	"github.com/michurin/cnbot/pkg/tests/apiserver"
	"github.com/michurin/cnbot/pkg/xbot"
	"github.com/michurin/cnbot/pkg/xctrl"
	"github.com/michurin/cnbot/pkg/xloop"
	"github.com/michurin/cnbot/pkg/xproc"
)

// *****************************************************************************
// Integration tests use real curl utility. It must be installed in your system.
// Motivation: to cover real behavior of curls command line options.
// *****************************************************************************

// TODO setup xlog

func TestAPI_justCall(t *testing.T) {
	/* case
	tg        bot
	|         |
	|<--req---|
	|---resp->|
	*/
	tgURL, tgClose := apiserver.APIServer(t, nil, map[string][]apiserver.APIAct{
		"/botMORN/xMorn": {{
			IsJSON:   true,
			Request:  `{"ok":1}`,
			Response: []byte(`{"response":1}`),
		}},
	})
	defer tgClose()

	ctx := context.Background()

	bot := buildBot(tgURL)

	body, err := bot.API(ctx, &xbot.Request{
		Method:      "xMorn",
		ContentType: "application/json",
		Body:        []byte(`{"ok":1}`),
	})

	require.NoError(t, err)
	assert.JSONEq(t, `{"response":1}`, string(body))
}

func TestMethods(t *testing.T) {
	/* cases
	tg        bots loop
	|         |
	|<--req---| (call for update)
	|---resp->|
	|         |
	|         |--exec-->| script
	|         |<-stdout-|
	|         |
	|<--req---| (call for update)
	|---resp->|
	|<--req---| (and send response from script)
	|---resp->| (the order of update and send doesn't meter)
	*/
	for _, cs := range []struct {
		name           string
		updateResponse string
		sendRequest    string
	}{
		{
			name: "message",
			updateResponse: `{"ok": true, "result": [{"update_id": 500, "message": {
"message_id": 100,
"from": {"id": 1500, "is_bot": false, "first_name": "Alex", "last_name": "Morn", "username": "AlexMorn", "language_code": "en"},
"chat": {"id": 1501, "first_name": "Alex", "last_name": "Morn", "username": "AlexMorn", "type": "private"},
"date": 1682222222,
"text": "word"}}]}`,
			sendRequest: `{"chat_id": 1500, "text": "word [n=1]"}`,
		},
		{
			name: "message_reaction",
			updateResponse: `{"ok": true, "result": [{"update_id": 500, "message_reaction": {
"message_id": 100,
"user": {"id": 1500, "is_bot": false, "first_name": "Alex", "last_name": "Morn", "username": "AlexMorn", "language_code": "en"},
"chat": {"id": 1501, "first_name": "Alex", "last_name": "Morn", "username": "AlexMorn", "type": "private"},
"date": 1682222222,
"old_reaction": [],
"new_reaction": [{"type":"emoji","emoji":"\ud83e\udd1d"}]}}]}`,
			sendRequest: `{"chat_id": 1500, "text": "message_reaction [n=1]"}`,
		},
		{
			name: "callback_query",
			updateResponse: `{"ok": true, "result": [{"update_id": 500, "callback_query": {
"id": "333333333333333333",
"from": {"id": 1500, "is_bot": false, "first_name": "Alex", "last_name": "Morn", "username": "AlexMorn", "language_code": "en"},
"chat": {"id": 1501, "first_name": "Alex", "last_name": "Morn", "username": "AlexMorn", "type": "private"},
"message": {
 "message_id": 90,
 "from": {"id": 1600, "is_bot": true, "first_name": "BOT", "username":"BOT_bot"},
 "date": 1682222222,
 "text": "OK",
 "reply_markup": {"inline_keyboard": [[{"text": "button_text", "callback_data": "button_data (in message)"}]]}},
"chat_instance": "4444444444444444444",
"data": "button_data"}}]}`,
			sendRequest: `{"chat_id": 1500, "text": "button_data [n=1]"}`,
		},
	} {
		cs := cs
		t.Run(cs.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tgURL, tgClose := apiserver.APIServer(t, cancel, map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": {
					{
						IsJSON:   true,
						Request:  `{"offset":0,"timeout":30,"allowed_updates":["callback_query","inline_query","message","message_reaction","poll","poll_answer"]}`,
						Response: []byte(cs.updateResponse),
					},
					{
						IsJSON:   true,
						Request:  `{"offset":501,"timeout":30,"allowed_updates":["callback_query","inline_query","message","message_reaction","poll","poll_answer"]}`,
						Response: nil,
					},
				},
				"/botMORN/sendMessage": {
					{
						IsJSON:   true,
						Request:  cs.sendRequest,
						Response: []byte(`{"ok": true, "result": {}}`),
					},
				},
			})
			defer tgClose()

			bot := buildBot(tgURL)

			command := buildCommand(t, "scripts/show_args.sh")

			err := xloop.Loop(ctx, bot, command)
			require.Error(t, err)
			require.Contains(t, err.Error(), "context canceled") // like "api: client: Post \"http://127.0.0.1:34241/botMORN/getUpdates\": context canceled"
		})
	}
}

func TestScriptOutputTypes(t *testing.T) { //nolint:funlen
	/* cases
	tg        bots loop
	|         |
	|<--req---| (call for update)
	|---resp->|
	|         |
	|         |--exec-->| script
	|         |<-stdout-|
	|         |
	|<--req---| (call for update)
	|---resp->|
	|<--req---| (and send response from script)
	|---resp->| (the order of update and send doesn't meter)
	*/
	simpleUpdates := []apiserver.APIAct{
		{
			IsJSON:  true,
			Request: `{"offset":0,"timeout":30,"allowed_updates":["callback_query","inline_query","message","message_reaction","poll","poll_answer"]}`,
			Response: []byte(`{"ok": true, "result": [{"update_id": 500, "message": {
"message_id": 100,
"from": {"id": 1500, "is_bot": false, "first_name": "Alex", "last_name": "Morn", "username": "AlexMorn", "language_code": "en"},
"chat": {"id": 1501, "first_name": "Alex", "last_name": "Morn", "username": "AlexMorn", "type": "private"},
"date": 1682222222,
"text": "word"}}]}`),
		},
		{
			IsJSON:   true,
			Request:  `{"offset":501,"timeout":30,"allowed_updates":["callback_query","inline_query","message","message_reaction","poll","poll_answer"]}`,
			Response: nil, // the second update call will stop Mock API server by this nil
		},
	}
	sendMessageResponseJSON := []byte(`{"ok": true, "result": {}}`)
	for _, cs := range []struct {
		script string
		api    map[string][]apiserver.APIAct
	}{
		{
			script: "scripts/text_ok.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendMessage": {
					{
						IsJSON:   true,
						Request:  `{"chat_id": 1500, "text": "ok"}`,
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			script: "scripts/text_long.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendMessage": {
					{
						IsJSON:   true,
						Request:  `{"chat_id": 1500, "text": "` + strings.Repeat(`⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘\n`, 315) + `1` + `"}`,
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			script: "scripts/text_too_long.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendDocument": {
					{
						IsJSON:   false,
						Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"document\"; filename=\"message.txt\"\r\nContent-Type: text/plain\r\n\r\n" + strings.Repeat("⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘\n", 315) + `12` + "\r\n--BOUND--\r\n",
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			script: "scripts/preformatted_ok.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendMessage": {
					{
						IsJSON:   true,
						Request:  `{"chat_id": 1500, "text": "ok", "entities": [{"type": "pre", "offset": 0, "length": 2}]}`,
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			script: "scripts/preformatted_complex_ok.sh", // one unicode char, however it is two utf16 words and length=2
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendMessage": {
					{
						IsJSON:   true,
						Request:  `{"chat_id": 1500, "text": "⚒️", "entities": [{"type": "pre", "offset": 0, "length": 2}]}`,
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			script: "scripts/preformatted_long.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendMessage": {
					{
						IsJSON:   true,
						Request:  `{"chat_id": 1500, "entities": [{"length":4096, "offset":0, "type":"pre"}], "text": "` + strings.Repeat(`⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘\n`, 315) + `1` + `"}`,
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			script: "scripts/preformatted_too_long.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendDocument": {
					{
						IsJSON:   false,
						Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"document\"; filename=\"formatted_text.txt\"\r\nContent-Type: text/plain\r\n\r\n" + strings.Repeat("⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘\n", 315) + `12` + "\r\n--BOUND--\r\n",
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			script: "scripts/media_jpeg.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendPhoto": {
					{
						IsJSON:   false,
						Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"photo\"; filename=\"image.jpeg\"\r\nContent-Type: image/jpeg\r\n\r\n\xff\xd8\xff\r\n--BOUND--\r\n",
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			script: "scripts/media_png.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendPhoto": {
					{
						IsJSON:   false,
						Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"photo\"; filename=\"image.png\"\r\nContent-Type: image/png\r\n\r\n\x89PNG\r\n\x1a\n\r\n--BOUND--\r\n",
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			script: "scripts/media_mp3.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendAudio": {
					{
						IsJSON:   false,
						Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"audio\"; filename=\"audio.mp3\"\r\nContent-Type: audio/mpeg\r\n\r\nID3\r\n--BOUND--\r\n",
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			script: "scripts/media_ogg.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendDocument": { // consider ogg as document, it seems it's not fully supported
					{
						IsJSON:   false,
						Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"document\"; filename=\"document.ogx\"\r\nContent-Type: application/ogg\r\n\r\nOggS\x00\r\n--BOUND--\r\n",
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			script: "scripts/media_mp4.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendVideo": {
					{
						IsJSON:   false,
						Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"video\"; filename=\"video.mp4\"\r\nContent-Type: video/mp4\r\n\r\n\x00\x00\x00\fftypmp4_\r\n--BOUND--\r\n",
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			script: "scripts/media_pdf.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendDocument": {
					{
						IsJSON:   false,
						Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"document\"; filename=\"document.pdf\"\r\nContent-Type: application/pdf\r\n\r\n%PDF-\r\n--BOUND--\r\n",
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			script: "scripts/media_bin.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendDocument": {
					{
						IsJSON:   false,
						Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"document\"; filename=\"document.dat\"\r\nContent-Type: application/octet-stream\r\n\r\n\x00\x00\x00\x00\r\n--BOUND--\r\n",
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			script: "scripts/media_len_zero.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
			},
		},
	} {
		cs := cs
		t.Run(cs.script[8:], func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tgURL, tgClose := apiserver.APIServer(t, cancel, cs.api)
			defer tgClose()

			bot := buildBot(tgURL)

			command := buildCommand(t, cs.script)

			err := xloop.Loop(ctx, bot, command)
			require.Error(t, err)
			require.ErrorContains(t, err, "context canceled") // like "api: client: Post \"http://127.0.0.1:34241/botMORN/getUpdates\": context canceled"
		})
	}
}

func TestHttp(t *testing.T) {
	/* cases
	tg        bots ctrl
	|         |
	|         |<-- someone external calls bot over http
	|<--req---| (request to send)
	|---resp->|
	|         |--> reply to external client
	*/
	for _, cs := range []struct {
		name string
		curl []string
		qs   string
		api  map[string][]apiserver.APIAct
	}{
		{
			name: "curl_F", // curl -F works transparently as is
			curl: []string{"-q", "-s", "-F", "user_id=10", "-F", "text=ok"},
			qs:   "",
			api: map[string][]apiserver.APIAct{
				"/botMORN/someMethod": {{
					IsJSON:   false,
					Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"user_id\"\r\n\r\n10\r\n--BOUND\r\nContent-Disposition: form-data; name=\"text\"\r\n\r\nok\r\n--BOUND--\r\n",
					Response: []byte("done."),
				}},
			},
		},
		{
			name: "curl_d",
			curl: []string{"-q", "-s", "-d", "ok"},
			qs:   "?to=111",
			api: map[string][]apiserver.APIAct{
				"/botMORN/sendMessage": {{
					IsJSON:   true,
					Request:  `{"chat_id":111, "text":"ok"}`,
					Response: []byte("done."),
				}},
			},
		},
	} {
		cs := cs
		t.Run(cs.name, func(t *testing.T) {
			tgURL, tgClose := apiserver.APIServer(t, nil, cs.api)
			defer tgClose()

			bot := buildBot(tgURL)

			h := xctrl.Handler(bot, nil, ctxlog.Patch(context.Background())) // we won't use second argument in this test

			s := httptest.NewServer(h)

			ou, er := runCurl(t, append(cs.curl, s.URL+"/x/someMethod"+cs.qs)...)
			assert.Equal(t, "done.", ou)
			assert.Empty(t, er)
		})
	}
}

func TestDownload(t *testing.T) {
	/* cases
	tg        bots ctrl
	|         |
	|         |<-- someone external calls bot over http with file_id
	|<--req---| (getFile)
	|---resp->|
	|<--req---| (download data)
	|---resp->|
	|         |--> reply to external client
	*/
	tgURL, tgClose := apiserver.APIServer(t, nil, map[string][]apiserver.APIAct{
		"/botMORN/getFile": {{
			IsJSON:   true,
			Request:  `{"file_id":"FILE"}`,
			Response: []byte(`{"ok":true, "result":{"file_path":"file/path.jpeg"}}`),
		}},
		"/file/botMORN/file/path.jpeg": {{
			IsJSON:   false,
			Stream:   true,
			Request:  "",
			Response: []byte("DATA"),
		}},
	})
	defer tgClose()

	bot := buildBot(tgURL)

	h := xctrl.Handler(bot, nil, ctxlog.Patch(context.Background())) // we won't use second argument in this test

	s := httptest.NewServer(h)

	ou, er := runCurl(t, "-q", "-s", s.URL+"?file_id=FILE")
	assert.Equal(t, "DATA", ou)
	assert.Empty(t, er)
}

func TestHttp_long(t *testing.T) { // CAUTION: test has sleep
	/* cases
	tg        bots loop
	|         |
	|         |<-- someone external calls bot over http (method=RUN)
	|         |
	|         |--exec-->| long-running external script
	|         |<-stdout-|
	|         |
	|<--req---| (request to send)
	|---resp->| (response will be skipped; and test tries cover it by making small sleep)
	*/

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tgURL, tgClose := apiserver.APIServer(t, cancel, map[string][]apiserver.APIAct{
		"/botMORN/sendMessage": {{
			IsJSON:   true,
			Request:  `{"chat_id":222, "text":"args [222]: a1 a2"}`,
			Response: nil, // response will be skipped, but in fact, we do not test this fact
		}},
	})
	defer tgClose()

	bot := buildBot(tgURL)
	command := buildCommand(t, "scripts/longrunning.sh")

	h := xctrl.Handler(bot, command, ctxlog.Patch(context.Background()))

	s := httptest.NewServer(h)

	ou, er := runCurl(t, "-q", "-s", "-X", "RUN", s.URL+"/?to=222&a=a1&a=a2")
	assert.Empty(t, ou)
	assert.Empty(t, er)
	<-ctx.Done()
	time.Sleep(time.Millisecond * 100) // we give small amount of time to let Bot.API method finishing after receiving response; it is not necessary
}

func TestProc(t *testing.T) { // CAUTION: test has sleep indirectly
	ctx := context.Background()
	t.Run("show_args", func(t *testing.T) {
		data, err := buildCommand(t, "scripts/run_show_args.sh").Run(ctx, []string{"ARG1", "ARG2"}, []string{"test1=TEST1", "test2=TEST2"})
		require.NoError(t, err, "data="+string(data))
		assert.Equal(t, "arg1=ARG1 arg2=ARG2 test1=TEST1 test2=TEST2 TEST=test\n", string(data))
	})
	t.Run("exit", func(t *testing.T) {
		data, err := buildCommand(t, "scripts/run_exit.sh").Run(ctx, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "wait: exit status 28")
		assert.Nil(t, data)
	})
	t.Run("slow", func(t *testing.T) {
		data, err := buildCommand(t, "scripts/run_slow.sh").Run(ctx, nil, nil)
		require.NoError(t, err)
		assert.Equal(t,
			`start
trap SIGINT
trap ERR
end
trap EXIT
`, string(data))
	})
	t.Run("immortal", func(t *testing.T) {
		data, err := buildCommand(t, "scripts/run_immortal.sh").Run(ctx, nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wait: signal: killed")
		assert.Nil(t, data)
	})
	t.Run("notfound", func(t *testing.T) {
		data, err := buildCommand(t, "scripts/NOTFOUND").Run(ctx, nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "/NOTFOUND: no such file or directory") // it is absolute path appears in error message
		assert.Nil(t, data)
	})
}

func buildCommand(t *testing.T, cmd string) *xproc.Cmd {
	t.Helper()
	absCmd, err := filepath.Abs(cmd) // app does it
	require.NoError(t, err)
	return &xproc.Cmd{
		InterruptDelay: 200 * time.Millisecond, // timeouts important for TestProc
		KillDelay:      200 * time.Millisecond,
		Env:            []string{"TEST=test"},
		Command:        absCmd,
	}
}

func buildBot(origin string) *xbot.Bot {
	return &xbot.Bot{
		APIOrigin: origin,
		Token:     "MORN",
		Client:    http.DefaultClient,
	}
}

func runCurl(t *testing.T, args ...string) (string, string) {
	t.Helper()
	t.Logf("Run curl %s", strings.Join(args, " "))
	cmd := exec.CommandContext(t.Context(), "curl", args...)
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	err := cmd.Run()
	require.NoError(t, err)
	return stdOut.String(), stdErr.String()
}
