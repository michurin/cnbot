#!/bin/sh

cmd="x_$(
    echo "$1" |
    tr '[:upper:][:space:]' '[:lower:]_' |
    sed 's-/--;s-___*-_-g; s-^_--; s-_$--'
)"
user="$BOT_TARGET_ID"
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
        convert rose: -resize 200x png:-
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
        echo '*Memory usage of PID='"$BOT_PID"'*'
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
    x_err_long)
        for i in `seq 4000`
        do
           echo $i
        done
        ;;
    x_err_invalid)
        printf '\xff\xff'
        ;;
    x_help)
        (
        echo '*Simplest examples:*'
        echo '`ok` — just say "ok"'
        echo '`nothig` — say nothing'
        echo '`empty` — empty message (no output)'
        echo '`rose` — send image'
        echo '`long` — emulate long running task; you can test, how concurrency works'
        echo '*Information:*'
        echo '`date` — date'
        echo '`pwd` — pwd'
        echo '`env` — env'
        echo '`mem` — memory usage'
        echo '`err` — emulate error exit'
        echo '`err_long` — emulate too long result string error'
        echo '`err_invalid` — emulate invalid UTF8 string error'
        echo '`help` — this message'
        echo '*Advanced examples:*'
        echo '`rrose` — image with cuption'
        echo '`note` — wait 3 seconds and push ordinary message'
        echo '`nnote` — wait 3 seconds and push message without notification'
        echo '`delayed` — delayed action'
        echo '*Experimental (raw messages)*'
        echo '`kbd` — keyboard'
        echo '`ikbd` — inline keyboard'
        echo '`del` — delete message'
        echo '`edit` — how bot can edit its messages'
        echo '*Misc. Just for fun*'
        echo '`cn` — random Chuck Norris joke from http://www.icndb.com'
        echo '`geo` — keyboard with requests of additional information'
        echo '`menu` — this help as menu'
        echo '`removemenu` — remove menu from screen'
        ) |
        curl -qsX POST -o /dev/null --data-binary @- "$url?parse_mode=markdown"
        echo .
        ;;
    x_rrose)
        convert rose: -resize 200x png:- |
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
    x_ikbd)
        cat <<JSON
sendMessage
{
    "chat_id": "$user",
    "text": "Demo of inline keyboard",
    "reply_markup": {"inline_keyboard": [
        [
            {"text": "google", "url": "http://google.com/"},
            {"text": "youtube", "url": "http://youtube.com/"}
        ], [
            {"text": "Say A", "callback_data": "A"},
            {"text": "Say B", "callback_data": "B"},
            {"text": "Say C", "callback_data": "C"}
        ]
    ]}
}
JSON
        ;;
    x_kbd)
        cat <<JSON
sendMessage
{
    "chat_id": "$user",
    "text": "Demo of keyboard",
    "reply_markup": {"keyboard": [
        [
            {"text": "help"},
            {"text": "date"},
            {"text": "env"}
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
    x_del)
        message_id=`echo 'Message will be deleted!' | curl -qsX POST --data-binary @- "$url"`
        sleep 2
        echo 'deleteMessage{"chat_id":'"$user"',"message_id":'"$message_id"'}' |
        curl -qsX POST -o /dev/null --data-binary @- "$url"
        echo .
        ;;
    x_edit)
        message_id=`echo 'Message will be edited!' | curl -qsX POST --data-binary @- "$url"`
        for text in 'Edited!' 'Updated text' 'Final text'
        do
            sleep 2
            method=editMessageText
            payload='{"chat_id":'"$user"',"message_id":'"$message_id"',"text":"'"$text"'"}'
            echo "$method$payload" |
            curl -qsX POST -o /dev/null --data-binary @- "$url"
        done
        echo .
        ;;
    x_cn)
        curl -s http://api.icndb.com/jokes/random |
        python -c 'import sys, json; t=json.load(sys.stdin); print(t["value"]["joke"])'
        ;;
    x_geo)
        cat <<JSON
sendMessage
{
    "chat_id": "$user",
    "text": "Demo of geo location and contact requests (not yet fully supported)",
    "reply_markup": {"keyboard": [
        [
            {"text": "Send location", "request_location": true},
            {"text": "Send contact", "request_contact": true}
        ]
    ]}
}
JSON
        ;;
    x_menu)
        cat <<JSON
sendMessage
{
    "chat_id": "$user",
    "text": "Help as menu, just push the button",
    "reply_markup": {"keyboard": [
JSON
        cat "$0" |
        perl -lne '
            BEGIN{$n=0; $a=$b=""}
            while(<>){if (m/^\s+x_([a-z_]+)\)\s*$/) {
                if ($n==0) {print "${a}["; $a="],"; $b=""}
                print qq[${b}{"text":"$1"}];
                $b=",";
                $n=($n+1)%4;
            }}
            END{print "]"}'
        cat <<JSON
    ]}
}
JSON
        ;;
    x_removemenu)
        cat <<JSON
sendMessage
{
    "chat_id": "$user",
    "text": "Remove keyboard",
    "reply_markup": {"remove_keyboard": true}
}
JSON
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
