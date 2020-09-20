#!/bin/sh

set -e

PATH=/usr/local/bin:/bin:/usr/bin

case "CMD_$1" in # NOTCMD
    CMD_sup)
        echo 'Hi there'
        ;;
    CMD_nothing)
        echo '.' # Single dot is marker of silence. The bot will reply nothing.
        ;;
    CMD_date)
        date
        ;;
    CMD_args)
        echo 'Passed args:'
        for a in "$@"
        do
            echo "=> $a"
        done
        ;;
    CMD_env)
        echo '```' # Marker of Markdown preformated message
        env
        echo '```'
        ;;
    CMD_rose)
        convert rose: -resize 200x png:- # Just throw image to stdout as is and it will appear in the chat
        ;;
    CMD_help)
        echo 'Available commands:' # We use "CMD"/"NOTCMD" substrings to build help message automatically
        grep CMD_ $0 | grep -v NOTCMD | sed 's-.*CMD_-=> -;s-.$--'
        ;;
    CMD_curl)
        url="http://localhost:9090/$BOT_NAME/to/$BOT_FROM"
        echo 'Message one' | curl -qsX POST -o /dev/null --data-binary @- "$url"
        echo 'Message two' | curl -qsX POST -o /dev/null --data-binary @- "$url"
        echo '.'
        ;;
    *)
        echo 'Try to say "help" to me'
        ;;
esac
