cnbot: Tool for creating Telegram bots as easy as writing sh script
===================================================================

What is it
----------

### Motivation

The goal of this project is to provide a way
to build Telegram bots as simple as CGI scripts.
All you need to write is a script (on any language)
that is complying with extremely simple contract.

### Key features

- Possibility to drive multiply bots by one process
- Extremely simple scripting thanks to concurrency policy and simple contract (see bellow)
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

#### [TODO] Setup bot

### [TODO] Run

#### [TODO] Chat with bot

#### [TODO] Send async messages

This bot can be used to send messages asynchronously.

### Possible problems

#### [TODO] Script not found

#### [TODO] Command not found

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
```
HTTPS_PROXY=socks5://localhost:8888 ./cnbot -c $PWD/config.json
```

Writing your own bot scenarios in details
------------------------------------------------

### [TODO] Arguments

### [TODO] Environment

### [TODO] Output treating

### [TODO] Life circle, timeouts and signals, concurrency and race conditions

### [TODO] Asynchronous flow

System administration topics
-----------------------------------

### [TODO] Daemonize and log rotation

### [TODO] Systemd example

Todo
----

- Modes of arguments preparation (for user's process)
- `reply_markup=InlineKeyboardMarkup` (?)
- Ability to send any Telegram API requests (using HTTP interface?)
- Travis
- Disable web hook if any `getWebhookInfo`+`deleteWebhook`

Relative links
--------------

- [Telegram Bot API](https://core.telegram.org/bots/api)
