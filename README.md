# cnbot: Tool for creating Telegram bots as easy as writing shell script

![build](https://github.com/michurin/cnbot/workflows/build/badge.svg)
![test](https://github.com/michurin/cnbot/workflows/test/badge.svg)
![lint](https://github.com/michurin/cnbot/workflows/lint/badge.svg)
[![@cnbot_demobot](http://shields.io/badge/demo_bot-%40cnbot__demobot-brightgreen?logo=telegram&style=flat)](https://t.me/cnbot_demobot)

The goal of this project is to provide a way
to alive Telegram bots by scripting that
even simpler than CGI scripts.
All you need to write is a script (on any language)
that is complying with extremely simple contract.

![Telegram bot screenshot](https://raw.githubusercontent.com/michurin/cnbot/static/screenshot.png)

## Features

### Quick glance

First of all, demo bot [@cnbot_demobot](https://t.me/cnbot_demobot)
is to give you an idea of how it works. This bot is driven
by `cnbot` and [this script](examples/public/script.sh).

You can find more details out of comments in
[demo script](examples/demo.sh) and [demo config](examples/config.json).
Read on "Quick start" section to learn how to run your own bot.

### Key features

- Extremely simple scripting
  - One request — one run (like CGI)
  - Using simply command line arguments, environment variables and `stdout`
  - You can just throw to `stdout` text or image, `cnbot` will distinguish what it is
  - `cnbot` keeps the queue of incoming messages and runs scripts strictly one after another. So you don't need to care about concurrent execution, locks and race conditions in your scripts
- Supports images
  - JPEG
  - PNG
  - GIF
- Supports
  - plain text messages
  - of preformatted text
  - markdown v2 formatted messages
  - inline keyboards, including mutable menus
- Supports asynchronous notification.
  Message could be emitted by the bot, not in reply only.
  You are free to send messages from `cron` scripts or
  other similar places
- Possibility to drive multiply bots by one process
- It's based on long polling. It means you don't have to have
  dedicated server, domain, SSL certificate, public IP or DMZ...
  You can run/develop/debug your bot straight from your laptop even behind NATs and firewalls
- Thanks to standard Go HTTP client, this bot supports proxy servers.
  So you can use it even if your provider have disallowed Telegram

### Basic values

- Simplicity of contract: you have to know only a few simple ideas to write your own bot
- Simplicity of code: functions better them methods, concrete types better than interfaces and so on
- Clear logs: caller, labels, clear and simple logging levels
- Security: restrict access by users white/black list,
  pass only certain environment variables to script,
  force working directory, timeouts,
  kill whole process group and so on

### Disadvantages and oversimplifications

- This bot is not designed for high load
- The bot doesn't have persistent storage.
  It can lose messages on a restart,
  no throttling, no retries
- This bot was not tested on MS Windows. It is further assumed that,
  you use Linux, macOS or other UNIX-like system
- This engine can control several bots driven by several
  scripts. However, there is no way to run those scripts
  with different permissions. It can be considered as misfeature

## Quick start

### Prerequisite

- You will need [golang](https://golang.org/) 1.14 or newer.
  You can [install](https://golang.org/doc/install) it even without root permissions, however,
  I will assume that you install it natively in your system.
- You have to [create](https://core.telegram.org/bots#3-how-do-i-create-a-bot)
  your Telegram bot and get your bot **token**.
  It is easy and totally free.
- It would be nice to know your Telegram user ID. You can
  learn it on the demo bot [@cnbot_demobot](https://t.me/cnbot_demobot)
  or special bot [@userinfobot](https://t.me/userinfobot).
  Or otherwise, you can proceed now, and find your ID out of
  the `cnbot` logs later.

### Compile

```sh
git clone https://github.com/michurin/cnbot.git
cd cnbot
go build ./cmd/...
ls -l cnbot
```

You'll get `cnbot` binary.

### Check configuration

```sh
export T='your Telegram API token here'
export U='your Telegram user ID, or just 0 in case you are going to get ID from cnbot logs'
./cnbot -i -c examples/config.yaml
```

You have to obtain your bot summary report.

### Run demo bot

Run it in the same way just without `-i` option.
Export `$T` and `$U` as above and run:

```sh
./cnbot -c examples/config.yaml
```

Now you can try to talk to your bot in your Telegram client.
The bot is set up to execute [examples/demo.sh](examples/demo.sh)
script.

### Writing your own bot

You can move forward using
[examples/config.yaml](examples/config.yaml)
[examples/demo.sh](examples/demo.sh).

Demo script demonstrate abilities of `cnbot` API.
You can play with it and read comment to figure out
how `cnbot` use arguments, environment variables and
treat scripts output.

Example of configuration file also contains all the
necessary guidance. It is recommended to *put your
token and user ID right into configuration file*
instead of environment variables as we did before.
It is more convenient and safe. You can alsa tune
timeouts, paths and few other options.

## Further reading

### Project documentation

#### How to start demo.sh-based bot and write your own scripts for bots

- All configuration options with comments: [configuration file example](examples/config.yaml)
- How to format messages and create inline keyboards: [example of script](examples/demo.sh)

#### Useful topics

- [Troubleshooting](doc/troubleshooting.md)
- [System administration topics](doc/system-administration.md): build for production, monitoring, startup scripts

Please feel free to email me at `a.michurin@gmail.com`

### References

- [Telegram Bot API](https://core.telegram.org/bots/api)
- Demo bot's icon made by [Smashicons](https://smashicons.com/) from [Flaticon](https://www.flaticon.com/)

## TODO

- Enrich `answerCallbackQuery` API call. Add corresponding types of control line
- Provide more user details through environment: last name, username, language.
- Add alive server information to environment
- Improve report (`-i` option)
- Improve project code layout
- Captions of images (using HTTP interface?)
- Modes of arguments preparation (for user's process)
- Ability to send any Telegram API requests (using HTTP interface?)
- Disable web hook if any `getWebhookInfo`+`deleteWebhook` (?)
- Play with `getMyCommands`/`setMyCommands` that were added in March 2020 in Bot API 4.7
- Examples of rc-scripts and systemd service-file.
- Test coverage
