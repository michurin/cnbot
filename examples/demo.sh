#!/bin/sh

set -e

PATH=/usr/local/bin:/bin:/usr/bin

case "CMD_$1" in # We are to say NOCMD here. See "help" section.
    CMD_sup)
        echo 'Hi there'
        ;;
    CMD_noout)
        ;;
    CMD_nothing)
        echo '.' # Single dot is marker of silence. The bot will reply nothing.
        ;;
    CMD_date)
        echo '```'
        date
        echo '```'
        ;;
    CMD_uname)
        echo '```'
        uname -a
        echo '```'
        ;;
    CMD_uptime)
        echo '```'
        uptime
        echo '```'
        ;;
    CMD_args)
        echo 'Passed args:'
        for a in "$@"
        do
            echo "- $a"
        done
        ;;
    CMD_env)
        echo '```' # Marker of Markdown preformatted message
        env | sort
        echo '```'
        ;;
    CMD_/env)
        u="http://$BOT_SERVER/$BOT_CHAT"
        (
        echo '```'
        env | sort
        echo '```'
        ) |
        curl -qfsX POST -o /dev/null --data-binary @- "$u"
        echo '.'
        ;;
    CMD_rose)
        montage rose: rose: rose: rose: -geometry +2+2 png:- # Just throw image to stdout as is and it will appear in the chat
        ;;
    CMD_du)
        d="$(df -P -m / | tail -1 | awk '{gsub("[^0-9]", "", $5); print $5","(100-$5)}')"
        # This old fashioned API is deprecated in 2012, however, it is still working
        # https://developers.google.com/chart/image/docs/making_charts
        u="https://chart.googleapis.com/chart?cht=p&chd=t:$d&chs=300x200&chl=Used|Available&chtt=Disk%20usage"
        curl -qfs "$u"
        ;;
    CMD_logo)
        curl -qfs https://golang.org/lib/godoc/images/footer-gopher.jpg
        # curl -qfs https://www.telegram.org/img/t_logo.png # try this if footer-gopher.jpg disappear
        ;;
    CMD_async)
        u="http://$BOT_SERVER/$BOT_FROM"
        echo "I'll send you logo..." | curl -qfsX POST -o /dev/null --data-binary @- "$u"
        curl -qfs https://www.telegram.org/img/t_logo.png | curl -qfsX POST -o /dev/null --data-binary @- "$u"
        echo 'Are you happy now?' | curl -qfsX POST -o /dev/null --data-binary @- "$u"
        echo '.'
        ;;
    CMD_help)
        echo 'Available commands:' # We use "CMD"/"NOCMD" substrings to build help message automatically
        grep CMD_ $0 | grep -v NOCMD | sed 's-.*CMD_-=> -;s-.$--'
        ;;
    *)
        echo 'Try to say "help" to me'
        ;;
esac
