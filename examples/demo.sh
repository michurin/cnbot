#!/bin/sh -e

# setup environment

PATH=/usr/local/bin:/bin:/usr/bin

# check if binaries exist

for c in awk cal curl date df echo env grep sed sort tail test uname uptime
do
  if ! command -v $c >/dev/null
  then
    echo "WARNING: Command $c not found, check PATH ($PATH)" >&2
  fi
done

# check if it is forwarded message or contact

if test -n "$BOT_SIDE_TYPE"
then
    echo "Message from $BOT_SIDE_TYPE"
    echo "Name: $BOT_SIDE_NAME"
    echo "ID: $BOT_SIDE_ID"
    exit
fi

# commands
# we use magic CMD marker to build help-message automatically

case "CMD_$1" in
    CMD_sup)
        echo 'Hi there! 👋'
        ;;
    CMD_noout)
        ;;
    CMD_nothing)
        echo '.' # Single dot is marker of silence. The bot will reply nothing.
        ;;
    CMD_date)
        date
        ;;
    CMD_uname)
        uname -a
        ;;
    CMD_uptime)
        uptime
        ;;
    CMD_args)
        echo 'Passed args:'
        for a in "$@"
        do
            echo "○ $a"
        done
        ;;
    CMD_env)
        echo '%!PRE'
        env | sort
        ;;
    CMD_calc)
        if test "$#" = '1'
        then
          echo '%!MARKDOWN'
          echo '*Usage:* calc _expression_'
          echo 'For example:'
          echo '`calc 1\+1`'
          echo '`calc atan2(1, 0) * 2`'
          echo '`calc sqrt(2)`'
        else
          prog="BEGIN {print($(echo "$BOT_TEXT" | sed 's/^[^c]*calc//'))}"
          echo '%!PRE'
          # echo "$prog" # uncomment for debugging
          awk "$prog" 2>&1
        fi
        ;;
    CMD_du)
        d="$(df -P -m / | tail -1 | awk '{gsub("[^0-9]", "", $5); print $5","(100-$5)}')"
        # This old fashioned API is deprecated in 2012, however, it is still working
        # https://developers.google.com/chart/image/docs/making_charts
        u="https://chart.googleapis.com/chart?cht=p&chd=t:$d&chs=300x200&chl=Available|Used&chtt=Disk%20usage"
        curl -qfs "$u"
        ;;
    CMD_gologo)
        curl -qfs https://golang.org/lib/godoc/images/footer-gopher.jpg
        # curl -qfs https://www.telegram.org/img/t_logo.png # try this if footer-gopher.jpg disappear
        ;;
    CMD_async)
        mark='%!MARKDOWN'$'\n' # it may surprise you, how we get newline in sh
        # the first way to send async message: multipart/form-data
        curl -qfsX POST -o /dev/null -F to=$BOT_FROM -F msg="${mark}_I'll send you_ *random* _image\.\.\._" $BOT_SERVER
        curl -qfsL https://source.unsplash.com/random/600x400 | curl -qfsX POST -o /dev/null -F to=$BOT_FROM -F msg=@- $BOT_SERVER
        # the second way to send async message: raw data + user_id at the end of url
        echo "${mark}_Are you happy now?_" | curl -qfsX POST -o /dev/null --data-binary @- "http://$BOT_SERVER/$BOT_FROM"
        echo '.' # suppress output processing
        ;;
    CMD_cal)
        echo '%!PRE'
        cal -h
        ;;
    CMD_btn)
        echo '%!CALLBACK env do env'
        echo '%!CALLBACK date do date'
        echo '%!CALLBACK'
        echo '%!CALLBACK uname do uname'
        echo '%!CALLBACK'
        echo '%!CALLBACK async get random image'
        echo '%!CALLBACK edit edit it!'
        echo 'Try buttons'
        ;;
    CMD_edit)
        echo '%!UPDATE'
        echo '%!CALLBACK env do env'
        echo '%!CALLBACK date do date'
        echo '%!CALLBACK'
        echo '%!CALLBACK uname do uname'
        echo '%!CALLBACK'
        echo '%!CALLBACK async get random image'
        echo '%!CALLBACK edit edit it!'
        #curl -qfsL 'http://api.forismatic.com/api/1.0/?method=getQuote&format=json&lang=en' | jq -r '"\(.quoteText) -- \(.quoteAuthor)"'
        curl -qfsL 'http://api.forismatic.com/api/1.0/?method=getQuote&format=text&lang=en'
        ;;
    CMD_help)
        echo '%!MARKDOWN'
        echo '*Available commands:*'
        grep CMD_ "$0" | grep -v case | sed 's-.*CMD_-• \\/-;s-.$--'
        echo ''
        echo '*And besides, the bot accespts*'
        echo '• contacts and'
        echo '• forwarded messages'
        echo 'to figure out user/chat/channel ID'
        ;;
    *)
        cmd="$(echo $1 | sed 's/[-_.]/\\&/g')" # we have to care about [-_.] only because other must-escaped-chars are disallowed by bot
        echo '%!MARKDOWN'
        echo "I didn't recognize your command '*$cmd*' Try to say '*help*' to me"
        ;;
esac
