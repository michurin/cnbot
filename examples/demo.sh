#!/bin/sh -e

PATH=/usr/local/bin:/bin:/usr/bin

for c in echo date uname uptime env sort curl tail awk grep sed cal
do
  if ! command -v $c >/dev/null
  then
    echo "Command $c not found, check PATH ($PATH)"
    exit
  fi
done

case "CMD_$1" in # We are to say NOCMD here. See "help" section.
    CMD_sup)
        echo 'Hi there! 👋'
        ;;
    CMD_noout)
        ;;
    CMD_nothing)
        echo '.' # Single dot is marker of silence. The bot will reply nothing.
        ;;
    CMD_date)
        echo '%!PRE'
        date
        ;;
    CMD_uname)
        echo '%!PRE'
        uname -a
        ;;
    CMD_uptime)
        echo '%!PRE'
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
          echo '%!MARKDOWN *Usage:* calc _expression_'
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
        u="https://chart.googleapis.com/chart?cht=p&chd=t:$d&chs=300x200&chl=Used|Available&chtt=Disk%20usage"
        curl -qfs "$u"
        ;;
    CMD_gologo)
        curl -qfs https://golang.org/lib/godoc/images/footer-gopher.jpg
        # curl -qfs https://www.telegram.org/img/t_logo.png # try this if footer-gopher.jpg disappear
        ;;
    CMD_async)
        u="http://$BOT_SERVER/$BOT_FROM"
        echo '%!MARKDOWN' "_I'll send you *random* image\.\.\._" | curl -qfsX POST -o /dev/null --data-binary @- "$u"
        curl -qfsL https://source.unsplash.com/random/600x400 | curl -qfsX POST -o /dev/null --data-binary @- "$u"
        echo '%!MARKDOWN' '_Are you happy now?_' | curl -qfsX POST -o /dev/null --data-binary @- "$u"
        echo '.'
        ;;
    CMD_cal)
        echo '%!PRE'
        cal -h
        ;;
    CMD_one)
        echo '%!MARKDOWN [go](tg://bot_command?command=one) \/one'
        ;;
    CMD_help)
        echo '%!MARKDOWN'
        echo '*Available commands:*' # We use "CMD"/"NOCMD" substrings to build help message automatically
        grep CMD_ $0 | grep -v NOCMD | sed 's-.*CMD_-• \\/-;s-.$--'
        ;;
    *)
        cmd="$(echo $1 | sed 's/[-.]/\\&/g')" # we have to care about [-.] only because other must-escaped-chars are disallowed by bot
        echo '%!MARKDOWN' "I didn't recognize your command '*$cmd*' Try to say '*help*' to me"
        ;;
esac
