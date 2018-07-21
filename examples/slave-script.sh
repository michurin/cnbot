#!/bin/sh

cmd="x_$(echo "$1" | tr '[:upper:][:space:]' '[:lower:]_' | sed 's-___*-_-g; s-^_--; s-_$--')"
user=$BOT_USER_ID
url="http://localhost:$BOT_SERVER_PORT/$user"

# TODO: check curl, convert, at

case "$cmd" in
    x_ok)
        echo ok
        ;;
    x_nothig)
        echo . # Single dot is a special silent marker, you will get nothing
        ;;
    x_empty)
        echo ''
        ;;
    x_rose)
        convert rose: png:-
        ;;
    x_long)
        steps=3
        for i in `seq $steps`
        do
            echo "OK PID=$$ step=$i/$steps" | curl -qsX POST --data-binary @- "$url" >/dev/null
            sleep 1
        done
        echo "OK PID=$$ FIN"
        ;;
    x_date)
        date
        ;;
    x_pwd)
        pwd
        ;;
    x_env)
        env | sort
        ;;
    x_mem)
        (
        echo '*Memory usage of PID=$BOT_PID*'
        echo '```'
        cat "/proc/$BOT_PID/status" | grep '^Vm' | expand
        echo '```'
        ) |
        curl -qsX POST -o /dev/null --data-binary @- "$url?parse_mode=markdown"
        echo .
        ;;
    x_err)
        echo '[stdout]'
        (>&2 echo "[stderr]")
        exit 1
        ;;
    x_help)
        (
        echo '*Simplest examples:*'
        echo '`ok` - just say "ok"'
        echo '`nothig` - say nothing'
        echo '`empty` - empty message (no output)'
        echo '`rose` - send image'
        echo '`long` - emulate long running task; you can test, how concurrency works'
        echo '*Information:*'
        echo '`date` - date'
        echo '`pwd` - pwd'
        echo '`env` - env'
        echo '`mem` - memory usage'
        echo '`err` - emulate error exit'
        echo '`help` - this message'
        echo '*Advanced examples:*'
        echo '`rrose` - image with cuption'
        echo '`note` - wait 3 seconds and push ordinary message'
        echo '`nnote` - wait 3 seconds and push message without notification'
        echo '`delayed` - delayed action'
        ) |
        curl -qsX POST -o /dev/null --data-binary @- "$url?parse_mode=markdown"
        echo .
        ;;
    x_rrose)
        convert rose: png:- |
        curl -qsX POST -o /dev/null --data-binary @- "$url?parse_mode=markdown&caption=Caption+text"
        echo .
        ;;
    x_note)
        sleep 3
        echo 'note! (with notification)' |
        curl -qsX POST -o /dev/null --data-binary @- "$url"
        echo .
        ;;
    x_nnote)
        sleep 3
        echo 'note! (without notification)' |
        curl -qsX POST -o /dev/null --data-binary @- "$url?disable_notification=yes"
        echo .
        ;;
    x_delayed)
        task_id=$$
        delayed_message="Task $task_id done."
        delayed_command="curl -qsX POST -o /dev/null --data-binary '$delayed_message' '$url'"
        echo "$delayed_command" | at -M now + 1minute >/dev/null 2>&1
        echo "Task $task_id scheduled"
        echo 'Delayed Command:'
        echo "$delayed_command"
        echo 'Wait one minute for result'
        ;;
    x_json)
        cat <<JSON
{
    "chat_id": "$user",
    "text": "Demo of inline keyboard",
    "reply_markup": {"inline_keyboard": [
        [
            {"text": "Open google.com", "url": "http://google.com/"},
            {"text": "Open youtube.com", "url": "http://youtube.com/"}
        ], [
            {"text": "Say A", "callback_data": "A"},
            {"text": "Say B", "callback_data": "B"},
            {"text": "Say C", "callback_data": "C"}
        ]
    ]}
}
JSON
        ;;
    x_callback_data:a)
        echo 'A!'
        ;;
    x_callback_data:b)
        echo 'B!'
        ;;
    x_callback_data:c)
        echo 'C!'
        ;;
    *)
        for i in "$@"
        do
          echo "ARG: [$i]"
        done
        echo "cmd: [$cmd]"
        echo "Say 'help' to get help"
        ;;
esac
