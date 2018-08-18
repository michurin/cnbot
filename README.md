cnbot
=====

Cnbot is successor to [instant-bot](https://github.com/michurin/instant-bot/).

Quick start
-----------

Edit `examples/config.toml` to use your bot token. Try to run:

```
dep ensure
go run cmd/bot/main.go -C examples/config.toml
```

Build for deploy and cross compilation
--------------------------------------

There three variables you can set in build time
to mark the version of binary: `Version`, `BuildRev` and `BuildDate`.

Example with cross compilation:

```
GOOS=linux GOARCH=386 go build -ldflags "-X main.Version=1 -X main.BuildRev=`git rev-parse HEAD` -X main.BuildDate=`date -u '+%Y-%m-%d_%H:%M:%S'`" -o bot-linux-386 cmd/bot/main.go
```

Full list of supported architectures you can find in
[Go sources](https://github.com/golang/go/blob/master/src/go/build/syslist.go).

On some platforms you may want to add
`--ldflags '-linkmode external -extldflags "-static"'`
to force build static binary.

Run
---

You can control logging using variables
`BOT_LOG_STDOUT`, `BOT_LOG_JSON` and `BOT_LOG_NODEBUG`.
For example:

```
BOT_LOG_STDOUT=1 BOT_LOG_NODEBUG=1 ./bot-linux-386
```

Also, you can use environment variable `HTTPS_PROXY=socks5://host:port`.
