#!/bin/bash

ALIVE_HANDLER='http://127.0.0.1:9999'

dirname="$(dirname $0)"
textdir="$dirname/texts"
datadir="$dirname/data"
logfile="$dirname/var/$(date +%y-%m-%d).log"
jq -c \
    --arg date "$(date +%s)" \
    --arg chat "$BOT_CHAT" \
    --arg from "$BOT_FROM" \
    --arg first_name "$BOT_FROM_FIRST_NAME" \
    --arg text "$BOT_TEXT" \
    --arg side_type "$BOT_SIDE_TYPE" \
    --arg side_name "$BOT_SIDE_NAME" \
    --arg side_id "$BOT_SIDE_ID" \
    --arg last_name "$BOT_FROM_LAST_NAME" \
    --arg username "$BOT_FROM_USERNAME" \
    --arg is_bot "$BOT_FROM_IS_BOT" \
    --arg language "$BOT_FROM_LANGUAGE" \
    --arg args "$*" \
    '. | .date=$date | .chat=$chat | .from=$from | .first_name=$first_name | .text=$text | .side_type=$side_type | .side_name=$side_name | .side_id=$side_id | .last_name=$last_name | .username=$username | .is_bot=$is_bot | .language=$language | .args=$args' <<<'{}' >>"$logfile"

if test -n "$BOT_SIDE_TYPE"
then
    echo 'It is forwarded message or contact'
    echo ''
    echo "Message from: $BOT_SIDE_TYPE"
    echo "Name: $BOT_SIDE_NAME"
    echo "ID: $BOT_SIDE_ID"
    exit
fi

case "CMD_$1" in
    CMD_date) # 💡 Execute command `date`
        date
        ;;
    CMD_uname) # 💡 `uname -a`
        uname -a
        ;;
    CMD_uptime) # 💡 `uptime`
        uptime
        ;;
    CMD_cal) # 💡 `cal`
        echo '%!PRE'
        cal -h
        ;;
    CMD_env) # 🤔 `env` \(show all environment variables available for script\) Try to say _"env Hello World\!"_
        echo '%!PRE'
        env | sort
        ;;
    CMD_args) # 🤔 Show scripts arguments\. Try to say _"args Hello World\!"_
        echo 'Passed args:'
        for a in "$@"
        do
            echo "👉 $a"
        done
        ;;
    CMD_image) # 🤪 Show random image
        curl -qfskL https://source.unsplash.com/random/600x400
        ;;
    CMD_buttons) # 🌶️ Buttons
        echo '%!CALLBACK date'
        echo '%!CALLBACK uname'
        echo '%!CALLBACK help list of commands (help)'
        echo '%!CALLBACK'
        echo '%!CALLBACK uptime Well, what uptime is it now? 🕓'
        echo 'Play with buttons!'
        ;;
    CMD_menu) # 🌶️ Mutable menu
        # According the rules of parsing, all invalid characters are considered
        # as separators. It means that "menu=1" turns to arg1=menu, arg2=1.
        # It is a way to make multi-arg callback.
        pg="$(sed 's/[^0-9]//g' <<< "$2")"
        pg="${pg:-0}"
        pg="$(( $pg % 5 ))"
        pga="$(( ($pg + 4) % 5 ))"
        pgb="$(( ($pg + 1) % 5 ))"
        echo '%!UPDATE'
        echo "%!CALLBACK menu=$pga << page $pga"
        echo "%!CALLBACK menu=$pgb page $pgb >>"
        echo '%!MARKDOWN'
        cat "$textdir/page$pg.txt" | sed 's/[-.(),]/\\&/g'
        ;;
    CMD_delayed) # 🌶️ Schedule delayed message \(background jobs with asynchronous response\)
        # Here we save "task" for cron.sh in a "quirky" way.
        # Our abuse protection may be considered as naive, however, it
        # works properly! Do not forget that cnbot cares about concurrency.
        # It executes scripts in a well-sequenced manner.
        if test "$(ls -1 "$datadir/"*-$BOT_CHAT-* | wc -l)" -gt 3
        then
            echo '%!MARKDOWN'
            echo 'Wait a moment\.\.\. Let me finish with your previous tasks\.'
        elif test "$(ls -1 "$datadir" | wc -l)" -gt 1000
        then
            echo '%!MARKDOWN'
            echo 'Sorry\. The bot is too busy\. Try later\.'
        else
            a=$RANDOM
            b=$RANDOM
            c=$RANDOM
            touch "$datadir/$(date "+%s-$BOT_CHAT-$a-$b-$c")" # We use file system as database.
            echo '%!MARKDOWN'
            echo "I promise you to solve the __*$a \\* $b \\* $c*__ for you in *two* minutes and come back with solution\\."
            echo 'In the meantime, we can just chat\.'
        fi
        ;;
    CMD_help) # ℹ️ List of commands
        echo '%!MARKDOWN'
        echo '*Available commands*'
        echo ''
        grep CMD_ "$0" | grep -v case | grep '#' | sed 's-.*CMD_-• \\/-;s-..#- —-'
        echo ''
        echo '*And besides, the bot accespts*'
        echo '• contacts and'
        echo '• forwarded messages'
        echo 'to figure out user/chat/channel ID'
        ;;
    CMD_info) # ℹ️ Build summary
        echo 'Build information:'
        echo ''
        curl -s "$ALIVE_HANDLER" | jq -r '"Build: \(.version.build)\nBot ver: \(.version.version)\nGo ver: \(.version.go)\nStarted at: \(.started_at)"'
        ;;
    CMD_start) # ℹ️ Show greeting
        cat "$textdir/start.txt"
        ;;
    *)
        cmd="$(sed 's/[-_.]/\\&/g' <<< "$1")"
        echo '%!MARKDOWN'
        echo "I didn't recognize your command '*$cmd*' Try to say '*help*' to me"
        ;;
esac

