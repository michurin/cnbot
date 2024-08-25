# cnbot

The goal of this project is to provide a way
to alive Telegram bots by scripting that
even simpler than CGI scripts.
All you need to write is a script (on any language)
that is complying with extremely simple contract.

![Telegram bot demo screenshot](https://raw.githubusercontent.com/michurin/cnbot/static/screenshot-2024.gif)

> [!NOTE]
> The minimal-effort echo-bot can be started like this:
>
> ````sh
> tb_token='TG_API_TOKEN' tb_script=echo tb_long_running_script=true tb_ctrl_addr=:9999 cnbot
> ````
>
> Where `echo` and `true` are standard command line utilities.

## What is it for

This bot engine has proven itself in alerting, system monitoring and managing tasks.

It also good for prototyping and fast proofing ideas.

## How mature is it

The engine is not perfect. Some error messages could be more informative.
Somewhere you can face a lug of documentation and the need to appeal to source code.

However, the engine has already proven itself in production and prototyping.

It served bots for huge conferences, meetings and events. It has helped customers
and provided control functionality for crew.

The engine successfully drives several monitoring and alerting bots.

It seems, API of this bot engines is quite stable and won't change dramatically in the near future.

## Basic ideas

You impalement all your business logic in your scripts. You are totally free to use all Telegram API abilities.

`cnbot` interact with scripts using (i) `stdout` stream, (ii) arguments and (iii) environment variables.

The engine automatically recognize multimedia and images. It cares about concurrency and races.

It also provides simple API for asynchronous messaging from `cron`s and such things.

It manages tasks (subprocesses), controls timeouts, sends signals and provides abilities to
run long-running tasks like long image/video conversions and/or downloading.

One instance of engine is able to manage several different bots.

## Quick start

### Run simplest one-line bot

#### Prepare

First things first, you need to create bot and get it's token.
It is free, just follow [instructions](https://core.telegram.org/bots#how-do-i-create-a-bot).

#### Build

TODO: `go install final_path`, hint: `GOBIN=$(pwd)`

#### Run

Just run one command to invoke the simplest bot:

```sh
tb_token='TOKEN' tb_script=/usr/bin/echo tb_long_running_script=/usr/bin/echo tb_ctrl_addr=:9999 cnbot
```

You are free to keep your token in file and use syntax like this to refer to file: `tb_token=@filename`

Don't worry, we will use configuration file further. The engine is able to use both files and direct environment variables.

- `tb_YOURBOTNAME_token` is a token your are given: `digits:long_string`
- `tb_YOURBOTNAME_script` is a command to run. We use the standard system command `echo`. I can be located elsewhere in your system. Try to say `whereis echo` to fine it
- `tb_YOURBOTNAME_long_running_script` let it be the same command. We consider it later
- `tb_YOURBOTNAME_ctrl_addr` we consider it soon

Run this command with correct variables and try to say something to you bot. You will be echoed by it.

### Put your configuration to file

You may as well put your configuration into env-file. The format of file is literally the same as `systemd` use.
So you are able to load it in `systemd` files as well. For example:

```sh
# let's name it config.env
tb_token='TOKEN'
tb_script=/usr/bin/echo
tb_long_running_script=/usr/bin/echo
tb_ctrl_addr=:9999
```

Now just start bot like this:

```sh
cnbot config.env
```

## Playing with random features

### Your first script (finding out your UserID)

Let's look at the script, that shows its arguments and environment variables:

```sh
#!/bin/sh

echo "Args: $@"
echo "Environment:"
env | grep tg_ | sort
```

Name it `mybot.sh` and mention it in configuration variable `tb_script=./mybot.sh`. Restart the bot and say to it `Hello bot!`.
It will reply to you something like that:

```
╭─────────────────────────────────────────╮
│ Args: hello bot!                        │
│ Environment:                            │
│ tg_message_chat_first_name=Alexey       │
│ tg_message_chat_id=153333328            │
│ tg_message_chat_last_name=Michurin      │
│ tg_message_chat_type=private            │
│ tg_message_chat_username=AlexeyMichurin │
│ tg_message_date=1717171717              │
│ tg_message_from_first_name=Alexey       │
│ tg_message_from_id=153333328            │
│ tg_message_from_is_bot=false            │
│ tg_message_from_language_code=en        │
│ tg_message_from_last_name=Michurin      │
│ tg_message_from_username=AlexeyMichurin │
│ tg_message_message_id=4554              │
│ tg_message_text=Hello bot!              │
│ tg_update_id=513333387                  │
│ tg_x_build=development (devel)          │
│ tg_x_ctrl_addr=:9999                    │
╰─────────────────────────────────────────╯
```

You can see that your message has been put to arguments in convenient normalized form, and you have a bunch of useful variables
with additional information. We will consider them further. At this point we just figure out then our user id is `tg_message_from_id=153333328`.
We will use this information very soon.

### Asynchronous messaging

You are free to send messages from anywhere: from cron jobs, from init scripts... Try it just from command line:

```sh
curl -qs http://localhost:9999/?to=153333328 -d 'OK!'
```

If you bot is running, you will obtain the message `OK!` in you Telegram client.

```
╭──────────╮
│ OK!      │
╰──────────╯
```

Do not forget to use *your* user id from previous section.

It makes sense what variable `tb_ctrl_addr=:9999` is for. It defines a control interface for external interactions with bot engine.

### Call arbitrary Telegram API methods

You can call whatever method you want. Full list of methods can be found in the
[official Telegram bot API documentation](https://core.telegram.org/bots/api).

For example, you can obtain information about your bot
(using method [getMe](https://core.telegram.org/bots/api#getme)):

```sh
curl -qs http://localhost:9999/method/getMe | jq
```

The response will look like this:

```json
{
  "ok": true,
  "result": {
    "id": 223333386,
    "is_bot": true,
    "first_name": "Your Bot",
    "username": "your_bot",
    "can_join_groups": true,
    "can_read_all_group_messages": false,
    "supports_inline_queries": false,
    "can_connect_to_business": false
  }
}
```

It enables you to send extended messages. For example, you can send a message with buttons
(method [sendMessage](https://core.telegram.org/bots/api#sendmessage)):

```sh
curl -qs http://localhost:9999/sendMessage -F chat_id=153333328 -F text='Select search engine' -F reply_markup='{"inline_keyboard":[[{"text":"Google","url":"https://www.google.com/"}, {"text":"DuckDuckGo","url":"https://duckduckgo.com/"}]]}'
```

You will receive message with two clickable buttons:

```
╭───────────────────────────╮
│ Select search engine      │
├─────────────┬─────────────┤
│ Google     ↗│ DuckDuckGo ↗│
╰─────────────┴─────────────╯
```

Do not forget to change `user_id`.

> [!NOTE]
> You can use any prefixes in URLs.
> URLs `http://localhost:9999/sendMessage` and `http://localhost:9999/ANITHING/sendMessage` are equal.
> It allows you to put engine's API behind prefix.

### Sending images

Bot recognizes media type of input. It will send text:

```sh
echo 'Hello!' | curl -qs http://localhost:9999/?to=153333328 --data-binary '@-'
```

However, it will send you image:

```sh
curl -qs https://github.githubassets.com/favicons/favicon.png | curl -qs http://localhost:9999/?to=153333328 --data-binary '@-'
```

> [!IMPORTANT]
> Please use the `--data-binary` option for binary data. Option `-d` corrupts EOLs.

### Formatted text

```sh
(echo '%!PRE'; echo 'Hello!') | curl -qs http://localhost:9999/?to=153812628 --data-binary '@-'
```

## Big picture

### Prepare playground

Let's extend our `mybot.sh` like that (it is literally [demo script](demo/demo_bot.sh) you can run by [docker compose](demo/compose.yaml)):

```sh
#!/bin/bash

LOG=logs/log.log # /dev/null

FROM="$tg_message_from_id"

API() {
    API_STDOUT "$@" >>"$LOG"
}

API_STDOUT() {
    url="http://localhost$tg_x_ctrl_addr/$1"
    shift
    echo "====== curl $url $@" >>"$LOG"
    curl -qs "$url" "$@" 2>>"$LOG"
    echo >>"$LOG"
    echo >>"$LOG"
}

(
    echo '==================='
    echo "Args: $@"
    echo "Environment:"
    env | grep tg_ | sort
    echo '...................'
) >>"$LOG"

case "$1" in
    debug)
        echo '%!PRE'
        echo "Args: $@"
        echo "Environment:"
        env | grep tg_ | sort
        echo "FROM=$FROM"
        echo "LOG=$LOG"
        ;;
    about)
        echo '%!PRE'
        API_STDOUT getMe | jq
        ;;
    two)
        API "?to=$FROM" -d 'OK ONE!'
        API "?to=$FROM" -d 'OK TWO!!'
        echo 'OK NATIVE'
        ;;
    buttons)
        bGoogle='{"text":"Google","url":"https://www.google.com/"}'
        bDuck='{"text":"DuckDuckGo","url":"https://duckduckgo.com/"}'
        API sendMessage \
            -F chat_id=$FROM \
            -F text='Select search engine' \
            -F reply_markup='{"inline_keyboard":[['"$bGoogle,$bDuck"']]}'
        ;;
    image)
        curl -qs https://github.com/fluidicon.png
        ;;
    invert)
        wm=0
        fid=''
        for x in $tg_message_photo # finding the biggest image but ignoring too big ones
        do
            v=${x}_file_size
            s=${!v} # trick: getting variable name from variable; we need bash for it
            if test $s -gt 102400; then continue; fi # skipping too big files
            v=${x}_width
            w=${!v}
            v=${x}_file_id
            f=${!v}
            if test $w -gt $wm; then wm=$w; fid=$f; fi
        done
        if test -n "$fid"
        then
            API_STDOUT '' -G --data-urlencode "file_id=$fid" -o - | mogrify -flip -flop -format png -
        else
            echo "attache not found (maybe it was skipped due to enormous size)"
        fi
        ;;
    reaction)
        API setMessageReaction \
            -F chat_id=$FROM \
            -F message_id=$tg_message_message_id \
            -F reaction='[{"type":"emoji","emoji":"👾"}]'
        echo 'Bot reacted to your message☝️'
        ;;
    madrid)
        API sendLocation \
            -F chat_id="$FROM" \
            -F latitude='40.423467' \
            -F longitude='-3.712184'
        ;;
    menu)
        mShowEnv='{"text":"show environment","callback_data":"menu-debug"}'
        mShowNotification='{"text":"show notification","callback_data":"menu-notification"}'
        mShowAlert='{"text":"show alert","callback_data":"menu-alert"}'
        mLikeIt='{"text":"like it","callback_data":"menu-like"}'
        mUnlikeIt='{"text":"unlike it","callback_data":"menu-unlike"}'
        mDelete='{"text":"delete this message","callback_data":"menu-delete"}'
        mLayout="[[$mShowEnv],[$mShowAlert,$mShowNotification],[$mLikeIt,$mUnlikeIt],[$mDelete]]"
        API sendMessage \
            -F chat_id=$FROM \
            -F text='Actions' \
            -F reply_markup='{"inline_keyboard":'"$mLayout"'}'
        ;;
    run)
        API "?to=$FROM&a=reactions&a=$tg_message_message_id" -X RUN
        echo "I'll show you long run"
        ;;
    edit)
        API "?to=$FROM&a=editing" -X RUN
        ;;
    id)
        echo '%!PRE'
        id 2>&1
        ;;
    caps)
        echo '%!PRE'
        getpcaps --verbose --iab $$
        ;;
    hostname)
        echo '%!PRE'
        hostname 2>&1
        ;;
    help)
        API sendMessage -F chat_id=$FROM -F parse_mode=Markdown -F text='
Known commands:

- `debug` — show args, environment and vars
- `about` — reslut of getMe
- `two` — one request, two responses
- `buttons` — message with buttons
- `image` — show image
- `invert` (as capture to image) — returns flipped flopped image
- `reaction` — show reaction
- `madrid` — show location
- `menu` — scripted buttons
- `run` — long-run example (long sequence of reactions)
- `edit` — long-run example (editing)
- `id` — check user who script runs from
- `caps` — check current capabilities (`getpcaps $$`)
- `hostname` — check hostname where script runs
- `help` — show this message
- `privacy` — mandatory privacy information
- `start` — just very first greeting message
'
        ;;
    start)
        API sendMessage -F chat_id=$FROM -F parse_mode=Markdown -F text='
Hi there!👋
It is demo bot to show an example of usage [cnbot](https://github.com/michurin/cnbot) bot engine.
You can use `help` command to see all available commands.'
        ;;
    privacy) # https://telegram.org/tos/bot-developers#4-privacy
        echo "This bot does not collect or share any personal information."
        ;;
    *)
        if test -n "$tg_callback_query_data"
        then
            case "$1" in
                menu-debug)
                    API answerCallbackQuery -F callback_query_id="$tg_callback_query_id"
                    echo '%!PRE'
                    echo "Environment:"
                    env | grep tg_ | sort
                    ;;
                menu-like)
                    API answerCallbackQuery -F callback_query_id="$tg_callback_query_id" -F "text=Like it"
                    API setMessageReaction -F chat_id=$tg_callback_query_message_chat_id \
                        -F message_id=$tg_callback_query_message_message_id \
                        -F reaction='[{"type":"emoji","emoji":"👾"}]'
                    ;;
                menu-unlike)
                    API answerCallbackQuery -F callback_query_id="$tg_callback_query_id" -F "text=Don't like it"
                    API setMessageReaction -F chat_id=$tg_callback_query_message_chat_id \
                        -F message_id=$tg_callback_query_message_message_id \
                        -F reaction='[]'
                    ;;
                menu-delete)
                    API answerCallbackQuery -F callback_query_id="$tg_callback_query_id"
                    API deleteMessage -F chat_id=$tg_callback_query_message_chat_id \
                        -F message_id=$tg_callback_query_message_message_id
                    ;;
                menu-notification)
                    API answerCallbackQuery -F callback_query_id="$tg_callback_query_id" -F text="Notification text (200 chars maximum)"
                    ;;
                menu-alert)
                    API answerCallbackQuery -F callback_query_id="$tg_callback_query_id" -F text="Notification text shown as alert" -F show_alert=true
                    ;;
            esac
        else
            API sendMessage -F chat_id=$FROM -F text='Invalid command. Say `help`.' -F parse_mode=Markdown
        fi
        ;;
esac
```

Let's add script for long-running tasks `mybot_long.sh` (it's [demo script](demo/demo_bot_long.sh)):

```sh
#!/bin/sh

LOG=logs/log_long.log # /dev/null

FROM="$tg_x_to"

API() {
    API_STDOUT "$@" >>"$LOG"
}

API_STDOUT() {
    url="http://localhost$tg_x_ctrl_addr/$1"
    shift
    echo "====== curl $url $@" >>"$LOG"
    curl -qs "$url" "$@" 2>>"$LOG"
    echo >>"$LOG"
    echo >>"$LOG"
}

case "$1" in
    reactions)
        MESSAGE_ID="$2"
        for e in "👾" "🤔" "😎"
        do
            API setMessageReaction -F chat_id=$FROM -F message_id=$MESSAGE_ID -F reaction='[{"type":"emoji","emoji":"'"$e"'"}]'
            sleep 1
        done
        API setMessageReaction -F chat_id=$FROM -F message_id=$MESSAGE_ID -F reaction='[]'
        ;;
    editing)
        MESSAGE_ID="$(API_STDOUT sendMessage -F chat_id=$FROM -F text='Starting...' | jq .result.message_id)"
        if test -n "$MESSAGE_ID"
        then
            for i in 2 4 6 8
            do
                sleep 1
                API editMessageText -F chat_id=$FROM -F message_id="$MESSAGE_ID" -F text="Doing... ${i}0% complete..."
            done
            sleep 1
            API editMessageText -F chat_id=$FROM -F message_id="$MESSAGE_ID" -F text='Done.'
        else
            echo "cannot obtain message id"
        fi
        ;;
    *)
        echo 'invalid mode'
        ;;
esac
```

Restart bot with this configuration (`mybot.env`):

```ini
tb_token               = 'TOKEN'
tb_script              = ./mybot.sh
tb_long_running_script = ./mybot_long.sh
tb_ctrl_addr           = :9999
```

Like that:

```sh
# if you install it
cnbot mybot.env
# if you start it without installing, just from sources
go run ./cmd/cnbot/... mybot.env
```

> [!NOTE]
> Please note when you are modifying script, all changes takes effect immediately. You don't need to restart the bot engine.
> You have to restart the bot engine if you want to change its environment variables only.

Try to talk to your bot. Now it recognizes commands and shows you many different possibilities.

Let me explain what is happening in this examples step by step.

### Script structure

You wouldn't be mistaken for thinking that this script is slightly awkward. It is written that way
to be more splittable. We will consider better structure further.

### Helpers overview

Let's briefly touch on two helpers functions we are using in this scripts.

Both of them helps you to call bot engine API (not Telegram API, but bot engine).

`API_STDOUT()` takes it's first argument as a tail of API URL and consider all the rest of arguments
as `curl`'s arguments. For example, `API_STDOUT getMe` means literally
`curl -qs "http://localhost$tg_x_ctrl_addr/getMe"`.

`API_STDOUT()` throws it's output to `stdout`, `API()` doesn't though.
`API "?to=$FROM" -d 'OK'` means `curl -qs "http://localhost$tg_x_ctrl_addr/?to=$FROM -d 'OK'`

Both of them logs their output to `$LOG` file.

### Commands

This script recognizes several commands. We already consider the following commands:

- `debug` — it's our first script
- `about` — just call `getMe` API method. You can also see how we use `API_STDOUT` helper
- `two` — shows how to send asynchronous message from script. We saw how to do it from command line before. You can also see how we use `API` helper
- `buttons` — message with buttons as we saw before
- `image` — shows how to send image. Just throw it to `stdout` and bot engine will recognize that it is image and send it in proper way

All the rest commands we will consider further.

## Advanced topics

### Configuration details and driving multiple bots

You are already seeing the bot can be configured by configuration file and directory by environment variable.

Environment has higher priority.

All variables have the same structure: `tb_{MEANING}` or `tb_{BOTNAME}_{MEANING}` if you need to start several bots.

To configure bot `x` and `y`, you need to pass this variable to `cnbot`:

```sh
tb_x_token='TOKEN_X'
tb_x_script=/usr/bin/echo
tb_x_long_running_script=/usr/bin/echo
tb_x_ctrl_addr=:9999

tb_y_token='TOKEN_Y'
tb_y_script=/usr/bin/echo
tb_y_long_running_script=/usr/bin/echo
tb_y_ctrl_addr=:9998
```

### Arguments processing

Bot engine runs your scripts with command line arguments. It can be useful for small bots.

Arguments prepared from messages, captions and callback's data. Strings are cast to lower-case, cleaned of control characters and split by white spaces.

For example the message `$Hello world!` will be represented as two arguments `hello` and `world`.

Following characters will be removed from the arguments: ``!"#$&'()*+-./:;<=>?@[\]`|``.

### Environment details

#### Turning telegram payload to environment variables

Bot engine converts every [JSON-update](https://core.telegram.org/bots/api#update) to flat set of environment variables this way:

```json
{
  "ok": true,
  "result": [
    {
      "message": {
        "caption": "Hi!",
        "chat": {
          "first_name": "Alexey",
          "id": 150000000,
          "last_name": "Michurin",
          "type": "private",
          "username": "AlexeyMichurin"
        },
        "date": 1600000000,
        "from": {
          "first_name": "Alexey",
          "id": 150000000,
          "is_bot": false,
          "language_code": "en",
          "last_name": "Michurin",
          "username": "AlexeyMichurin"
        },
        "message_id": 2222,
        "photo": [
          {
            "file_id": "aaa0",
            "file_size": 2444,
            "file_unique_id": "id0",
            "height": 90,
            "width": 90
          },
          {
            "file_id": "aaa1",
            "file_size": 4888,
            "file_unique_id": "id1",
            "height": 128,
            "width": 128
          }
        ]
      },
      "update_id": 500000000
    }
  ]
}
```

turns to the following environment variables:

```ini
tg_message_caption=Hi!
tg_message_chat_first_name=Alexey
tg_message_chat_id=150000000
tg_message_chat_last_name=Michurin
tg_message_chat_type=private
tg_message_chat_username=AlexeyMichurin
tg_message_date=1600000000
tg_message_from_first_name=Alexey
tg_message_from_id=150000000
tg_message_from_is_bot=false
tg_message_from_language_code=en
tg_message_from_last_name=Michurin
tg_message_from_username=AlexeyMichurin
tg_message_message_id=2222
tg_message_photo=tg_message_photo_0 tg_message_photo_1
tg_message_photo_0_file_id=aaa0
tg_message_photo_0_file_size=2444
tg_message_photo_0_file_unique_id=id0
tg_message_photo_0_height=90
tg_message_photo_0_width=90
tg_message_photo_1_file_id=aaa1
tg_message_photo_1_file_size=4888
tg_message_photo_1_file_unique_id=id1
tg_message_photo_1_height=128
tg_message_photo_1_width=128
tg_update_id=500000000
```

#### Build-in variables (`x`-variables)

Engine provides the following additional variables:

- `tg_x_build`
- `tg_x_ctrl_addr`
- `tg_x_to` (long-running scripts only)

#### System variables

> [!NOTE]
> Beware. Bot engine does *NOT* convey its environment to child scripts.

Bot engine does not transfer environment to child scripts. It is conscious decision cause it helps to
make script's behavior more predictable and reproducible. Variables like `$PATH`, `$LANG`, `$LS_ALL` can
change behavior of many commands and functions. It can lead to hard to debug behavior.

If you need to have some environment variables, just set them in you script explicitly.

### Working directory

Current working directory is directory, where the script is located in.

### Process management: concurrency, timeouts, signals, long-running tasks

#### Ordinary tasks

Bot engine generates all tasks of the same bot run strictly concurrently. It means you can use
shared resources like files without any doubts. And your tasks have to finish in short time.

Bot engine will send `SIGTERM` to task after 10 seconds, and `SIGKILL` after next 10 seconds.

#### Long-running tasks

Long-running tasks can be executed simultaneously though.

They also have timeouts: 10 minutes.

### Uploading and downloading

To upload something (image, video, audio, etc) you can just throw it stdout of your script.
If you need to add capture or group multimedia files in one message, you need to call
Telegram API. As usual, you don't need to care about secrets etc just use `cnbot` control handler as we did above.

To download attachments (file, video, audio, photos, etc) you have to use `file_id` from message and
just perform `GET` request to control handler with `file_id=...` in query string. See action `invert`
in example above.

## Tips and tricks

### Improved script structure and security aspects

```sh
# --- global variables
...
# --- helper variables
...
# --- must have commands
case $1 in
start)
    echo "Hello message"
    exit
    ;;
privacy) # https://telegram.org/tos/bot-developers#4-privacy
    echo "This bot does not collect or share any personal information."
    exit
esac
# --- whitelist checks for user_id
# it is just example:
# - allows.list have contains strings line "_${ID}_" (it makes you able to write comments and things like that)
# - we consider messages and callbacks
if grep "_${tg_message_from_id}${tg_callback_query_from_id}_" allows.list 2>&1 >/dev/null
then
    : # pass this user, you may want to log it
else
    echo 'You are not allowd'
    exit
fi
# --- process text messages
if [ -n "$tg_message_text" ]
then
    case "$1" in
        ...
    esac
    exit
fi
# --- process images
if [ -n "$tg_message_photo" ]
then
    case "$1" in
        ...
    esac
    exit
fi
# --- process voices (for instance)
if [ -n "$tg_message_voice_file_id" ]
then
    ...
    exit # don't forget to exit
fi
# --- process callbacks
if [ -n "$tg_callback_query_data" ]
then
    ...
    exit
fi
# process... whatever you want
if ...
    ...
    exit
fi
```

Of course, it is good idea to split script, using `source file.sh` instruction.
And you are still able to use other languages and approaches for sure.

> [!CAUTION]
> Just don't forget to be careful, keep in mind that anybody in internet can send anything to your bot.
>
> Keep reading. We will consider how to protect your bot.

### Debugging wrapper

To debug your scripts, you can use this wrapper. Tune `$CMD`, and enjoy
full logging: arguments, environment, out and err streams, exit code.

```sh
#!/bin/sh

# put your command here
CMD=./mybot.py

# tune naming for your taste
base="logs/$(date +%s-)_${$}_"
ext='.log'

n=0
for a in "$@"
do
    echo "$a" >"${base}arg_${n}${ext}"
    n="$(($n+1))"
done

env | sort >"${base}env${ext}"

set -o pipefail

"$CMD" "$@" 2>"${base}err${ext}" | tee "${base}out${ext}"

code="$?"

echo "$code" >"${base}status${ext}"
exit "$code"
```

## System administration topics

### Installation

```sh
./build.sh
sudo install ./cnbot /usr/bin
```

### Running

The process itself does not try to be immortal. It dies on fatal issues that can not be solved by process itself. Like network problems.
It is believed that the process will be restart by `systemd` or stuff like that according the proper way with timeouts, logging, notifications, alerting.

Systemd unit file example (`/etc/systemd/system/cnbot.service`):

```ini
[Unit]
Description=Telegram bot (cnbot) service
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
Restart=always
RestartSec=1
User=nobody
ExecStart=/usr/bin/cnbot /etc/cnbot-config.env

[Install]
WantedBy=multi-user.target
```

## Known issues

- Some engine API methods are using both POST-body and query parameters. It's against standards. However I haven't invented something more convenient and standard yet.
- Engine API uses non-standard method `RUN`. It allows by standards, however it doesn't seem inevitable.
- Engine doesn't retry any requests to Telegram API. Looks like issue. However, Telegram API doesn't provide any idempotency keys, and engine doesn't save state between restarts. It seems you have to solve this issue somehow else.
- It hasn't been tested on MS Windows and FreeBSD.
- The engine doesn't support persistent storage. You have to save state if you need by yourself.

## Developing and contributing

### Main ideas

- Contract must be simple and flexible
- New features of [Telegram bot API](https://core.telegram.org/bots/api) has to be available instantly without changing of code of the bot
- Bot has to manage subprocesses: timeouts, etc
- Bot has to manage API call: [rate limits](https://core.telegram.org/bots/faq#my-bot-is-hitting-limits-how-do-i-avoid-this), etc
- Configuration must be simple
- Code must be testable and has to be covered
- Functionality has to be observable and has to provide ability to add metrics and monitoring by adding middleware without code changing

### Deep debugging

Run proxy. For example [mitmproxy](https://mitmproxy.org/):

```sh
mitmdump --flow-detail 4 -p 9001 --mode reverse:https://api.telegram.org
```

Instruct the bot to use proxy and run it:

```sh
export tb_api_origin=http://localhost:9001
./cnbot ... # run bot, it will deal with Telegram API through the proxy and you will see everything
```

### Application structure

(horrible ASCII art warning)

```
   Telegram infrastructure
             ^                             ............. crons
        HTTP :                        HTTP :             scripts
             :                             v             any other
.=BOT================================================.   asynchronous
|            API           | HTTP server for         |
|..........................| asynchronous messaging  |
| polling for : sending    |                         |
| updates     : messages  <-- send data from req     |
`===================================================='
    |             ^    ^  send stdout     |
    |             |    `---------.        | request params
    | message     | send         |        | as command line positional args
    v data        | stdout       |        v
........................        ......................
: run script for every :        : long-running       :
: message              :        : script             :
:......................:        :....................:
```
