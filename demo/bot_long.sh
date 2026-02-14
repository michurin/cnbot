#!/bin/bash

LOG=logs/log_long.log # /dev/null to disable logging
mkdir -p "$(dirname "$LOG")" # do not forget to create all necessary directories

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
            API setMessageReaction -F chat_id="$FROM" -F message_id="$MESSAGE_ID" -F reaction='[{"type":"emoji","emoji":"'"$e"'"}]'
            sleep 1
        done
        API setMessageReaction -F chat_id="$FROM" -F message_id="$MESSAGE_ID" -F reaction='[]'
        ;;
    editing)
        MESSAGE_ID="$(API_STDOUT sendMessage -F chat_id="$FROM" -F text='Starting...' | jq .result.message_id)"
        if test -n "$MESSAGE_ID"
        then
            for i in 2 4 6 8
            do
                sleep 1
                API editMessageText -F chat_id="$FROM" -F message_id="$MESSAGE_ID" -F text="Doing... ${i}0% complete..."
            done
            sleep 1
            API editMessageText -F chat_id="$FROM" -F message_id="$MESSAGE_ID" -F text='Done.'
        else
            echo "cannot obtain message id"
        fi
        ;;
    progress)
        MESSAGE_ID="$(API_STDOUT sendMessage -F chat_id="$FROM" -F text='Starting...' | jq .result.message_id)"
        if test -n "$MESSAGE_ID"
        then
            for i in \
                '[..........]' \
                '[#.........]' \
                '[##........]' \
                '[###.......]' \
                '[####......]' \
                '[#####.....]' \
                '[######....]' \
                '[#######...]' \
                '[########..]' \
                '[#########.]'
            do
                sleep 1
                API editMessageText -F chat_id="$FROM" -F message_id="$MESSAGE_ID" -F text="\`${i}\` Doing..." -F parse_mode=Markdown
            done
            sleep 1
            API editMessageText -F chat_id="$FROM" -F message_id="$MESSAGE_ID" -F text='Done.'
        else
            echo "cannot obtain message id"
        fi
        ;;
    *)
        echo 'invalid mode'
        ;;
esac
