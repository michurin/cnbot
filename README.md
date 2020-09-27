cnbot: Tool for creating Telegram bots as easy as writing sh script
===================================================================

What is it
----------

### Motivation

The goal of this project is to provide a way
to alive Telegram bots by scripting that
even simpler than CGI scripts.
All you need to write is a script (on any language)
that is complying with extremely simple contract.

### Key features

- Possibility to drive multiply bots by one process
- Extremely simple scripting thanks to strict concurrency policy and simple contract (see bellow)
- Support of images (JPEG/PNG/GIF)
- Support of preformatted text (markdown)
- Support asynchronous notification (message could be emitted by the bot, not in reply only)
- It's based on polling. It means you don't have to have dedicated server, domain, public IP or DMZ... You can run/develop/debug your bot straight from your laptop
- Thanks to standard Go HTTP client, this bot supports proxy servers. So you can use it even if your provider disallows Telegram

### Basic values

- Simplicity: functions better them methods, concrete types better then interfaces
- Clear logs: caller, labels
- Security: restrict access by users white list, pass only certain environment variables to script

Quick start
-----------

### Compile

All you need is [Go](https://golang.org/) language 1.13
or newer. You can [install](https://golang.org/doc/install) it
even without root permissions. To compile just say:

```sh
git clone https://github.com/michurin/cnbot.git
cd cnbot
go build ./cmd/...
```

and you get `cnbot` binary.

### Configure

#### Get bot token

First of all you have to
[create](https://core.telegram.org/bots#3-how-do-i-create-a-bot)
your Telegram bot and get your bot **token**.

#### Setup bot

Create in the root of project the following file

```json
{
  "bots": [
    {
      "token": "111:xxx",
      "script": "examples/demo.sh",
      "allowed_users": [
        153812628
      ],
      "working_dir": "/tmp",
      "term_timeout": 3,
      "kill_timeout": 1,
      "wait_timeout": 1,
      "bind_address": ":9090",
      "read_timeout": 2,
      "write_timeout": 2
    }
  ]
}
```

Quick glance over it:

- `bots` Array of bots. One process is able to drive several bots.
    - `token` Your token. Usually looks like `[degits]:[garbage]`.
    - `script` Path to your executable script. Absolute or relative directory with configuration file.
    - `working_dir` Working dir. Absolute or relative directory with configuration file.
    - `allowed_users` White list of users. The easiest way to figure out your ID is to leave this array
    empty, try to talk to bot and find in logs error message about rejected user.
    - `term_timeout`, `kill_timeout`, `wait_timeout`. Optional. Delays in second. Could be fractional.
    - `bind_address` Optional. Address in format like `127.0.0.1:8000`, `[::]:8080` or just `:8888`.
    - `read_timeout`, `write_timeout` Optional. Request reading and writing timeouts in seconds. Could be fractional.

The details are given below.

### Run

```sh
./cnbot -c $PWD/config.json
```

TODO: screenshots

If you faced problems, read no.

### Chat with bot

TODO: screenshots

#### Send async messages

This bot can be used to send messages asynchronously.

Try this commad with your bot name and user ID.

```sh
date | curl -qfsvX POST -o /dev/null --data-binary @- 'http://localhost:9090/${YOUR_BOT_NAME}/to/${TARGET_USER_ID}'
```

You can send images and preformatted text as well.

```sh
cat YOUR_IMAGE.jpeg | curl -qfsvX POST -o /dev/null --data-binary @- 'http://localhost:9090/${YOUR_BOT_NAME}/to/${TARGET_USER_ID}'
```

See more details bellow and in `examples/demo.sh`.

### Possible problems

#### Script not found: no such file or directory

If you receive error `no such file or directory`,
please check the `"script"` attribute in
your configuration file. You can
specify an absolute path, or relative path
related to directory of configuration file.

#### Command not found

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

#### HTTP port is already in use

If you receive log message like that
```
listen tcp :9090: bind: address already in use
```
it means what it means. If you don't need Web server
for asynchronous messages, just remove `server` section
from config. Otherwise, change ports to eliminate
conflict.

#### Two bots at the same time

It is impossible to run more than one bot. Otherwise you
will receive `HTTP 409 Conflict` error with payload like this:
```json
{"ok":false,"error_code":409,"description":"Conflict: terminated by other getUpdates request; make sure that only one bot instance is running"}
```
You are to stop one bot.

#### Web hook

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

#### Telegram API is unreachable (use proxy)

Use `HTTPS_PROXY` environment
```sh
HTTPS_PROXY=socks5://localhost:8888 ./cnbot -c $PWD/config.json
```

Writing your own bot scenarios in details
-----------------------------------------

As you saw above, you are to write you
"script" to animate your bot. This script
can be written on any language and executed
bo every user message as separate process.
It's very similar to CGI-scripts approach.

### Arguments

Whole message is casted to lower case, split by spaces
and passed to script as a set of arguments.

For example, the message `Hello world! ` will be represented
by two arguments: `hello` and `world!`.

### Environment

Only three environment variables are passed
to script:

- `BOT_NAME` is a name of bot that receive a message.
It is useful if your script drive several bots simultaneously.
- `BOT_FROM` is the sender's user id. You can use it
to schedule asynchronous actions and
corresponding messaging.
- `BOT_SERVER` is a real bind address for asynchronous communication.

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

#### Preformatted reply

If your output stream is started and ended by
the triple backtick `` ``` `` sequence. All spaces
around this triple are ignored.

You don't need to care about escaping special chars
in your message.

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

If `"server"` key is present in bot's configuration, `cnbot` starts for that bot simple HTTP server.

It accepts the `POST` method and paths like `/${TARGET_USER}`. `cnbot` threats
the body of request in the same way as output of scripts. So you can send text, preformatted text
and images in PNG, JPEG and GIF formats.

Using this mechanism you can:

- Send asynchronous messages
- Send messages to another users
- Even send messages from another bots

You can configure reading and writing timeouts. Default values are 10s and 10s.

System administration topics
----------------------------

### Daemonize

`cnbot` doesn't have any daemonization abilities.
If you really wish to daemonize it, you can use tools like `nohup`.
However, if you use `sistemd` you don't need for daemonization.

### Log management

`cnbot` just throws log messages to `stdout`. You can use
tools like `multilog` to manage log files.

Todo
----

- Modes of arguments preparation (for user's process)
- `reply_markup=InlineKeyboardMarkup` (?)
- Ability to send any Telegram API requests (using HTTP interface?)
- Travis
- Disable web hook if any `getWebhookInfo`+`deleteWebhook`
- Play with `getMyCommands`/`setMyCommands` that were added in March 2020 in Bot API 4.7
- Tests
- Throttling
- Examples of rc-scripts and systemd service-file.

Relative links
--------------

- [Telegram Bot API](https://core.telegram.org/bots/api)
