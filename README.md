# cnbot

Cnbot is successor to [instant-bot](https://github.com/michurin/instant-bot/).

## Key idea

The goal of this project to provide a simple way to implement Telegram
bot with arbitrary functionality using any languages and tools.

You have just to save you script, configure `cnbot` and run it. When
bot receives a message, it runs your script and pass message to it,
your script have to reply by text message or even by *image*, and bot
send the reply back to Telegram user.

Moreover, you can send messages asynchronously (not in reply). Bot
provides HTTP interface to do it.

## Quick start

First of all you have to [create](https://core.telegram.org/bots#3-how-do-i-create-a-bot)
your Telegram bot and get your bot **token**.

Download bot sources and all dependencies; and build binary:
```sh
cd ~/tmp # anywhere you wish
export GOPATH=`pwd`
go get -v github.com/michurin/cnbot/...
cp src/github.com/michurin/cnbot/examples/* .
vim config.toml # chenge token to your token
# export HTTPS_PROXY=socks5://localhost:8888 # if you need to use proxy
bin/bot -C config.toml
```

Try to open Telegram application and chat with bot.

## Build for deploy and cross compilation

### Build with meta

There three variables you can set in build time
to mark the version of binary: `Version`, `BuildRev` and `BuildDate`.

Example with cross compilation:

```sh
GOOS=linux GOARCH=386 go build -ldflags "-X main.Version=1 -X main.BuildRev=`git rev-parse HEAD` -X main.BuildDate=`date -u '+%Y-%m-%d_%H:%M:%S'`" -o bot-linux-386 cmd/bot/main.go
```

Full list of supported architectures you can find in
[Go sources](https://github.com/golang/go/blob/master/src/go/build/syslist.go).

On some platforms you may want to add
`--ldflags '-linkmode external -extldflags "-static"'`
to force build static binary.

### Run

You can control logging using variables
`BOT_LOG_STDOUT`, `BOT_LOG_JSON` and `BOT_LOG_NODEBUG`.
For example:

```sh
BOT_LOG_STDOUT=1 BOT_LOG_NODEBUG=1 ./bot-linux-386
```

Also, you can use environment variable `HTTPS_PROXY=socks5://host:port`.

## Configuration and protocol

Under construction. See [examples](examples):
[configuration file](examples/config.toml) and
[script](examples/slave-script.sh) with bot commands implementation.

## License

Copyright (c) 2018, Alexey Michurin \<a.michurin@gmail.com\>. All rights reserved.

Licensed under the [New BSD (no advertising, 3 clause)](LICENSE) License.
