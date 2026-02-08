# Docker cheat sheet

## About container

This `Dockerfile` provides container for development, testing and playing around. Not for production.
It contains tools for editing and viewing, for sniffing network traffic, for debugging scripts and verbose logging.

## Prerequisites

All you need to run demo `cnbot` is Telegram bot token. You can [obtain it for free](https://core.telegram.org/bots/tutorial#obtain-your-bot-token).

## First run and first glance

Build image

```sh
sudo docker build -t cnbot:v2 .
```

> [!TIP]
> You can specify branch using `docker build --build-arg branch=experimental`.

Run bot (use your own token)

```sh
sudo docker run -it --rm --name cnbot -e 'TB_TOKEN=4839574812:AAFD39kkdpWt3ywyRZergyOLMaJhac60qc' cnbot:v2
```

Now bot is ready, you can to talk to it in your Telegram client.

## Enter to container, viewing logs and playing with scripts

You are free to enter to container

```sh
sudo docker exec -it cnbot /bin/bash
```

If you want to be `root` in the container, you can:

```sh
sudo docker exec -it -u root cnbot /bin/bash
```

Now you are in container. You can run `mc` for navigation, `vim` or `ne` for editing.

You can find logs in `/app/logs`, you can modify `/app/bot.sh` and `/app/bot_logn.sh`. You do not need to restart `cnbot` to see your changes.

You can turn on super verbose logging

```sh
sudo docker run -it --rm --name cnbot -e 'TB_TOKEN=4839574812:AAFD39kkdpWt3ywyRZergyOLMaJhac60qc' -e 'TB_SCRIPT=/app/bot_debug.sh' -e 'TB_LONG_RUNNING_SCRIPT=/app/bot_long_debug.sh' cnbot:v2
```

## Calling bot's control API in container

For example to call `getMe` just say

```sh
sudo docker exec cnbot curl -qs http://localhost:9999/getMe | jq
```

## Examination interaction with Telegram API

It is `mitmproxy` available in the container. You can run `mitmproxy` tool like this

```sh
sudo docker run -it --rm --name cnbot -e 'TB_TOKEN=4839574812:AAFD39kkdpWt3ywyRZergyOLMaJhac60qc' -e 'TB_SCRIPT=/app/bot_debug.sh' -e 'TB_LONG_RUNNING_SCRIPT=/app/bot_long_debug.sh' -e 'TB_API_ORIGIN=http://localhost:9001' cnbot:v2 /usr/local/bin/mitmdump --set confdir=/tmp --flow-detail 4 -p 9001 --mode reverse:https://api.telegram.org
```

And than execute `cnbot` in this container:

```sh
sudo docker exec -it cnbot ./cnbot
```

Now try to talk with bot and enjoy verbose output.

## Just docker memo

```sh
sudo docker images                 # list images
sudo docker image rm -f cnbot:v2   # remove image

sudo docker image prune            # cleanup everything
sudo docker container prune
```

You are free to install what you need in container

```sh
apt-get update
apt-get install -y zsh
```
