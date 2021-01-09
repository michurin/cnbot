#!/bin/bash

# You can find more explanations in demo.sh

export TZ=UTC

alivehandler='http://127.0.0.1:8900'

dirname="$(dirname $0)"
textdir="$dirname/texts"
datadir="$dirname/data"
logfile="$dirname/var/$(date +%y-%m-%d).log"

jq -n -c \
    --arg date "$(date +%s)" \
    --arg chat "$BOT_CHAT" \
    --arg from "$BOT_FROM" \
    --arg first_name "$BOT_FROM_FIRST_NAME" \
    --arg text "$BOT_TEXT" \
    --arg message_type "$BOT_MESSAGE_TYPE" \
    --arg side_type "$BOT_SIDE_TYPE" \
    --arg side_name "$BOT_SIDE_NAME" \
    --arg side_id "$BOT_SIDE_ID" \
    --arg last_name "$BOT_FROM_LAST_NAME" \
    --arg username "$BOT_FROM_USERNAME" \
    --arg is_bot "$BOT_FROM_IS_BOT" \
    --arg language "$BOT_FROM_LANGUAGE" \
    --arg loc_long "$BOT_LOCATION_LONGITUDE" \
    --arg loc_lat "$BOT_LOCATION_LATITUDE" \
    --arg loc_acc "$BOT_LOCATION_ACCURACY" \
    --arg loc_period "$BOT_LOCATION_LIVE_PERIOD" \
    --arg loc_heading "$BOT_LOCATION_HEADING" \
    --arg loc_alert_r "$BOT_LOCATION_ALERT_RADIUS" \
    --arg args "$*" \
    '{
    "date": $date,
    "chat": $chat,
    "from": $from,
    "first_name": $first_name,
    "text": $text,
    "message_type": $message_type,
    "side_type": $side_type,
    "side_name": $side_name,
    "side_id": $side_id,
    "last_name": $last_name,
    "username": $username,
    "is_bot": $is_bot,
    "language": $language,
    "loc_long": $loc_long,
    "loc_lat": $loc_lat,
    "loc_acc": $loc_acc,
    "loc_period": $loc_period,
    "loc_heading": $loc_heading,
    "loc_alert_r": $loc_alert_r,
    "args": $args
    }' >>"$logfile"

if test -n "$BOT_SIDE_TYPE"
then
    echo 'It is forwarded message or contact'
    echo ''
    echo "Message from: $BOT_SIDE_TYPE"
    echo "Name: $BOT_SIDE_NAME"
    echo "ID: $BOT_SIDE_ID"
    exit
fi

if test -n "$BOT_LOCATION_LONGITUDE" # BOT_LOCATION_LATITUDE can be used too
then
    echo '%!PRE'
    echo 'Location is received'
    echo ''
    echo "Longitude: $BOT_LOCATION_LONGITUDE"
    echo "Latitude:  $BOT_LOCATION_LATITUDE"
    echo ''
    echo 'Optional:'
    echo "Accuracy:     $BOT_LOCATION_ACCURACY"
    echo "Live period:  $BOT_LOCATION_LIVE_PERIOD"
    echo "Heading:      $BOT_LOCATION_HEADING"
    echo "Alert radius: $BOT_LOCATION_ALERT_RADIUS"
    exit
fi

case "${BOT_MESSAGE_TYPE}_$1" in
    message_date)
        date
        ;;
    message_uptime)
        uptime
        ;;
    message_cal)
        echo '%!PRE'
        cal -h
        ;;
    message_env)
        echo '%!PRE'
        # it has to work correctly with multiline values
        python -c 'import os; print("\n".join(sorted(k+"="+v for k, v, in os.environ.items() if k.startswith("BOT_"))))'
        ;;
    message_args)
        echo 'Passed args:'
        for a in "$@"
        do
            echo "👉 $a"
        done
        ;;
    message_image)
        curl -qfskL https://source.unsplash.com/random/600x400
        ;;
    message_buttons)
        echo '%!CALLBACK buttons-helper-date date as notification'
        echo '%!CALLBACK buttons-helper-uname uname as alert'
        echo '%!CALLBACK'
        echo '%!CALLBACK buttons-helper-uptime Well, what uptime is it now? 🕓'
        echo 'Play with buttons!'
        ;;
    callback_buttons-helper-date)
        echo '%!SILENT'
        echo "%!TEXT $(date)"
        ;;
    callback_buttons-helper-uname)
        echo '%!SILENT'
        echo "%!ALERT $(uname -a)"
        ;;
    callback_buttons-helper-uptime)
        uptime
        ;;
    message_menu|callback_menu)
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
    message_delayed)
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
    message_help)
        echo '%!MARKDOWN'
        cat "$textdir/help.md"
        ;;
    message_info)
        echo 'Build information:'
        echo ''
        curl -s "$alivehandler" | jq -r '"Build: \(.version.build)\nBot ver: \(.version.version)\nGo ver: \(.version.go)\nStarted at: \(.started_at)"'
        ;;
    message_start)
        echo '%!MARKDOWN'
        cat "$textdir/start.txt"
        ;;
    *)
        cmd="$(sed 's/[-_.]/\\&/g' <<< "$1")"
        echo '%!MARKDOWN'
        echo "I didn't recognize your __${BOT_MESSAGE_TYPE}__ '__${cmd}__' Try to say '*__help__*' to me"
        ;;
esac
