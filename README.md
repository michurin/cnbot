# cnbot: Tool for creating Telegram bots as easy as writing shell script

![build](https://github.com/michurin/cnbot/workflows/build/badge.svg)
![test](https://github.com/michurin/cnbot/workflows/test/badge.svg)
![lint](https://github.com/michurin/cnbot/workflows/lint/badge.svg)

## :mag_right: What it is

It is a tool to create custom simple Telegram bots. For example, this is
a part of chat with [demo script](examples/demo.sh).

![Telegram bot screenshot](https://raw.githubusercontent.com/michurin/cnbot/static/screenshot.png)

### Motivation

The goal of this project is to provide a way
to alive Telegram bots by scripting that
even simpler than CGI scripts.
All you need to write is a script (on any language)
that is complying with extremely simple contract.

### Key features

- Extremely simple scripting
  - One request - one run (like CGI)
  - Using simply command line arguments, environment variables and `stdout`
  - You can just throw to `stdout` text or image, `cnbot` will distinguish what it is
  - `cnbot` handles incoming messages strictly one after another. So you don't need to care about concurrency and race conditions in your scripts
- Possibility to drive multiply bots by one process
- Support of JPEG, PNG and GIF formats
- Support of preformatted text (markdown)
- Support asynchronous notification. Message could be emitted by the bot, not in reply only. You are free to send messages from `cron` scripts and other similar places
- It's based on polling. It means you don't have to have dedicated server, domain, SSL certificate, public IP or DMZ... You can run/develop/debug your bot straight from your laptop even behind NATs and firewalls
- Thanks to standard Go HTTP client, this bot supports proxy servers. So you can use it even if your provider have disallowed Telegram

### Basic values

- Simplicity of contract: you have know only a few simple ideas to write your own bot
- Simplicity of code: functions better them methods, concrete types better then interfaces and so on
- Clear logs: caller, labels, clear logging levels
- Security: restrict access by users white list, pass only certain environment variables to script, force working directory, timeouts, kill whole process group and so on

### Disadvantages and oversimplifications

- This bot is not design for high load
- The bot doesn't have persistent storage: it can lose messages on restart, no throttling, no retries
- Inline keyboard, custom reply keyboard, messages editing and other features are not supported

## :airplane: Quick start

### Compile

All you need is [Go](https://golang.org/) language 1.14
or newer. You can [install](https://golang.org/doc/install) it
even without root permissions. To compile just say:

```sh
git clone https://github.com/michurin/cnbot.git
cd cnbot
go build ./cmd/...
```

and you'll get `cnbot` binary.

### Configure

#### Get bot token

First of all you have to
[create](https://core.telegram.org/bots#3-how-do-i-create-a-bot)
your Telegram bot and get your bot **token**.

#### Setup bot

Create configuration file in the root of the project.

You can start from simplest one like this with your token and user id.

```yaml
bots:
  simplest_echo_bot:
    token: "111:xxx"
    script: /bin/echo
    allowed_users: [153812628]
```

You can find out your user id in several ways:

- You can use bot like [@userinfobot](https://t.me/userinfobot)
- You can just start cnbot with empty `allowed_users`, try to add the bot in your Telegram client, find out your user id from bot's error message `user [your id] is not allowed`

This bot just echoes your messages. However, let's consider full configuration
file and full-futured [demo script](examples/demo.sh):

```yaml
bots:
  bot_nickname:
    token: "111:xxx"
    script: "examples/demo.sh"
    allowed_users: [153812628]
    working_dir: "/tmp"
    term_timeout: 3
    kill_timeout: 1
    wait_timeout: 1
    bind_address: ":9090"
    read_timeout: 2
    write_timeout: 2
```

Quick glance over it:

- `bots` Set of bots. One process is able to drive several bots. Bots `nicknames` are for internal purpose only like logging.
    - `token` Your token. Usually looks like `[degits]:[garbage]`.
    - `script` Path to your executable script. Absolute or relative.
    - `working_dir` The working dir. Absolute or relative.
    - `allowed_users` White list of users. The easiest way to figure out your ID is to leave this array
    empty, try to talk to bot and find in logs error message about rejected user.
    - `term_timeout`, `kill_timeout`, `wait_timeout`. Optional. Delays in second. Could be fractional.
    - `bind_address` Optional. Address in format like `127.0.0.1:8000`, `[::]:8080` or just `:8888`.
    - `read_timeout`, `write_timeout` Optional. Request reading and writing timeouts in seconds. Could be fractional.

Relative paths in `script` and `working_dir` is treated with respect to the configuration file directory-path.

The details are given below.

### Check configuration file

```sh
./cnbot -i -c config.yaml
```

Note that the path to the configuration file must be absolute.

If your configuration file is ok, you will receive bot summary for every configured bot.
Report looks like that:

```
REPORT

- go version: go1.15.2 / linux / amd64
- nickname: "test1"
  - bot info:
    - bot id: 1036010710
    - bot name: "test_111_bot"
    - first name: "Test 111"
    - can join grp: true
    - can read all grp msgs: false
    - support inline: false
    - web hook: empty (it's ok)
  - configuration:
    - allowed users: 153812628
    - script:
      - script: "/home/a/cnbot/examples/demo.sh"
      - working dir: "/home/a/cnbot"
      - timeouts: 5s, 500ms, 500ms (term/kill/wait)
    - server:
      - address: ":9091"
      - timeouts: 10s, 10s (w/r)
```

### Run

```sh
./cnbot -c config.yaml
```

If you faced problems, read no.

### Chat with bot

#### Send async messages

This bot can be used to send messages asynchronously.

Try this command with your user ID.

```sh
date | curl -qfsvX POST -o /dev/null --data-binary @- 'http://localhost:9090/${TARGET_USER_ID}'
```

You can send images and preformatted text as well.

```sh
cat YOUR_IMAGE.jpeg | curl -qfsvX POST -o /dev/null --data-binary @- 'http://localhost:9090/${TARGET_USER_ID}'
```

See more details bellow and in [examples/demo.sh](examples/demo.sh).

## :pill: Troubleshooting

### Script not found: no such file or directory

If you receive error `no such file or directory`,
please check the `"script"` attribute in
your configuration file. You can
specify an absolute path, or relative path
related to directory of configuration file.

### Command not found

If your script use subprocesses, it is possible
it will unable to find corresponding binaries.

It is because `cnbot` doesn't pass its
environment to subprocesses.
Long story short, you can specify full
absolute paths or export `PATH` environment variable.

Why does `cnbot` behave is that way?
Where are two crucial reasons:

- Repeatable behavior. You script will word in the same environment on your laptop,
on the dedicated server, from  ordinary user and from root.
- Security reasons. You have to specify `PATH`
explicitly in your script.

Long story. If you use some kind of `shell`
for scripting you may be surprised to know
that `PATH` is a magic variable. It exists
even if it doesn't. At least it is true for
[bash](http://git.savannah.gnu.org/cgit/bash.git/tree/variables.c#n486)
([DEFAULT_PATH_VALUE](http://git.savannah.gnu.org/cgit/bash.git/tree/config-top.h#n66)
defined on compilation time) and
[zsh](https://github.com/zsh-users/zsh/blob/master/Src/init.c#L978).

You may want to set variables like
`LD_PRELOAD`, `LANG`, `LC_*`,
`POSIXLY_CORRECT`, `PYTHONPATH`, `PERL5LIB`
along with `PATH`.

### HTTP port is already in use

If you receive log message like that
```
listen tcp :9090: bind: address already in use
```
it means what it means. If you don't need Web server
for asynchronous messages, just remove `server` section
from config. Otherwise, change ports to eliminate
conflict.

### Two bots at the same time

It is impossible to run more than one bot. Otherwise you
will receive `HTTP 409 Conflict` error with payload like this:
```json
{"ok":false,"error_code":409,"description":"Conflict: terminated by other getUpdates request; make sure that only one bot instance is running"}
```
You are to stop one bot.

### Web hook

If your bot has been drown by another tool/library/SDK and
this bot starts to receive `HTTP 409 Conflict` error with
following payload:
```
{"ok":false,"error_code":409,"description":"Conflict: can't use getUpdates method while webhook is active; use deleteWebhook to delete the webhook first"}
```
It means you set up web hook. You can check it by command like that: 
```sh
curl 'https://api.telegram.org/bot$TOKEN/getWebhookInfo'
```
```json
{"ok":true,"result":{"url":"https://your.url/","has_custom_certificate":false,"pending_update_count":0,"max_connections":40,"allowed_updates":["message"]}}
```
You can easily remove web hook manually:
```sh
curl 'https://api.telegram.org/bot$TOKEN/deleteWebhook'
```
```json
{"ok":true,"result":true,"description":"Webhook was deleted"}
```

### Telegram API is unreachable (use proxy)

Use `HTTPS_PROXY` environment
```sh
HTTPS_PROXY=socks5://localhost:8888 ./cnbot -c config.yaml
```

## :wrench: Writing your own bot scenarios in details

As you saw above, you are to write you
"script" to animate your bot. This script
can be written on any language and executed
bo every user message as separate process.
It's very similar to CGI-scripts approach.

### Arguments

`cnbot`

- casts message to lower case
- split message by disallowed characters to arguments list. Allowed characters:
  - Letters: `a`–`z`
  - Digits: `0`–`9`
  - Dot, underscore, minus: `.`, `_`, `-`

For example, the message `Hello world! ` will be represented
by two arguments: `hello` and `world`.
Message `/go@bot run!` will turn to three arguments `go`, `bot` and `run`.

Original text of message is still available in `$BOT_TEXT`. See below.

### Environment

Only three environment variables are passed
to script:

- `BOT_NAME` — a name of bot that receive a message according configuration file
- `BOT_CHAT` — chat id. It is used to reply
- `BOT_FROM` — the sender's user id
- `BOT_FROM_FIRSTNAME` — sender's first name
- `BOT_SERVER` — a real bind address for asynchronous communication
- `BOT_TEXT` — raw message

There are additional variables for forwarded messages and contacts:

- `BOT_SIDE_TYPE` — type of data source: `user`, `bot`, `contact`, `channel`, `private`, `group` or `supergroup`
- `BOT_SIDE_ID` — user id or chat id
- `BOT_SIDE_NAME` — user name, bot name, channel title or contact name

Be aware that all other variables are not passed,
including `PATH`. See notes above.

### Exit code

Script have to finish with zero exit code. All
other values are treated as error.

### Output treating

`cnbot` can distinguish several types of output.

#### No output

If bot receive nothing from script, it sends to user the word _empty_ (italic).

It is convenient in many cases when you want to see output of commands like `who`,
that can produce nothing.

#### Suppressed reply

If you really don't want to reply to message, you need to produce a single dot `.`.
In this case bot won't send anything in reply.

#### Formatting

If your output stream starts by
the `%!PRE` signature, it will be escaped properly and represented as preformatted text

If your output stream starts by
the `%!MARKDOWN` signature, it will be sent as MarkdownV2. You *must* to escape all
special chars according [documentation](https://core.telegram.org/bots/api#markdownv2-style).

#### Images

Bot can detect three types of images: PNG, JPEG and PNG.
If they are detected, they send as proper images.

It is very convenient to send graphs from Graphana, RRD tools
or something like that.

### Stderr

All stderr output appears in logs only.

### Timeouts and signals

You can set up tree timeout: _term_, _kill_ and _wait_.

After _term_ timeout script is notified by `SIGTERM`.
After that, `cnbot` waits _kill_ timeout and send `SIGKILL`.
Next, `cnbot` waits _wait_ timeout for exit status.

All signals are directed to a process group, so all scripts children are signaled too.

Default values of timeouts are: 10s, 1s, 1s.

### Concurrency and race conditions

`cnbot` guarantees that scripts are never execute simultaneously.
You a free to use sheared resources without worry about locks, concurrency and races.

It is not scalable solution. You can change this behaviour easily, just split `msgQueue`.
However, be careful with that, remember about throttling. Currently, it doesn't
implement.

### Asynchronous flow

Due to security reasons, each bot has its own server.

If `bind_address` key is present in bot's configuration, `cnbot` starts a simple HTTP server for that bot.

The server threats
messages in the same way as output of scripts. So you can send text, preformatted text
and images in PNG, JPEG and GIF formats.

Using this mechanism you can:

- Send asynchronous messages
- Send messages to another users
- Even send messages from another bots

You can configure reading and writing timeouts. Default values are 10s and 10s.

There are two similar ways to send messages asynchronously.

#### `multipart/form-data` request with `to` and `msg` params (`POST`)

```sh
curl -F to=$BOT_FROM -F msg=text $BOT_SERVER
```

#### Raw `POST` with message in body and target user in URL

```sh
curl --data-binary text "http://$BOT_SERVER/$BOT_FROM"
```

See examples in [demo script](examples/demo.sh).

## :pizza: System administration topics

### Build

You can specify build version

```sh
go build -ldflags "-s -w -X github.com/michurin/cnbot/pkg/bot.Build=`date +%F`-`git rev-parse --short HEAD`" ./cmd/...
```

### Monitoring and health checking

You can turn on alive handler, mentioning it in configuration file like that:

```yaml
bot: ...
alive:
  bind_address: ":9999"
```

After that you can fetch current state of process from this handler: version, start time, memory utilisation details, current number of goroutines...
All this information is available in json format, so you can use tools like `jq` to manage it. For example, it is easy to make input string for
simple RRDtools based monitoring:

```sh
curl -s localhost:9999 | jq -r '"N:\(.num_goroutine):\(.mem_status.Alloc):\(.mem_status.HeapObjects)"'
```

It will return string like `N:16:1365296:10822` suitable for `rrdtool update`.

### Startup: modern `systemd` fashion

As usual create service file like this `/etc/systemd/system/cnbot.service`:

```ini
[Unit]
Description=Telegram bot (cnbot) service
After=network.target

[Service]
Type=simple
Restart=always
RestartSec=1
User=nobody
ExecStart=/usr/bin/cnbot -c /etc/cnbot-config.yaml

[Install]
WantedBy=multi-user.target
```

Don't forget to set your actual username after `User=` and the proper path to your script
and configuration file in `ExecStart=`. Keep in mind that relative paths in configuration file
evaluate with respect of configuration file directory.

```
[root]# systemctl daemon-reload
[root]# systemctl status cnbot
● cnbot.service - Telegram bot (cnbot) service
     Loaded: loaded (/etc/systemd/system/cnbot.service; disabled; vendor preset: disabled)
     Active: inactive (dead)
[root]# systemctl start cnbot
[root]# systemctl status cnbot # you will see that bot is active
```

If everything is ok, you can enable service (`systemd enable`).
If something goes wrong, you can inspect logs:

```sh
journalctl -u cnbot
```

You can use `-f` with `journalctl` to read live-tail logs.

### Startup: old fashioned `rc.d` scripts

#### Daemonize

`cnbot` doesn't have any daemonization abilities.
If you really wish to daemonize it, you can use tools like `nohup`.
However, if you use `sistemd` you don't need for daemonization.

#### Log management

`cnbot` just throws log messages to `stdout`.
There is a lot of tools designed with the primary purpose
of maintaining size-capped, automatically rotated, log file sets.
Personally, I prefer `multilog` from [daemontools](http://cr.yp.to/daemontools.html).

Further reading: [Don't use logrotate or newsyslog in this century](http://jdebp.eu/FGA/do-not-use-logrotate.html).

## :bulb: Todo

- Modes of arguments preparation (for user's process)
- `reply_markup=InlineKeyboardMarkup` (?)
- Ability to send any Telegram API requests (using HTTP interface?)
- Disable web hook if any `getWebhookInfo`+`deleteWebhook` (?)
- Play with `getMyCommands`/`setMyCommands` that were added in March 2020 in Bot API 4.7
- Tests
- Examples of rc-scripts and systemd service-file.

## :link: Relative links

- [Telegram Bot API](https://core.telegram.org/bots/api)
