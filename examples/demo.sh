#!/bin/sh -e

# This script is to show all features of cnbot.
#
# Generally, cnbot creates messages out of two sources:
# - it receives messages from Telegram users, runs script according configuration file and treats script's output as response
# - it listens HTTP requests (if configured), treats bodies of requests in the same way and sends the result to Telegram users
#
# So, this script is run for every user request and it's output is used to reply to user.

# PART I: Set up environment #########################################
#
# cnbot doesn't pass environment to script to avoid side effects and prevent possible vulnerabilities. So, you are to
# set up explicitly variables like PATH, LANG, and so on.

PATH=/usr/local/bin:/bin:/usr/bin

# PART II: Check if binaries exist ##################################
#
# You are free to skip it.

for c in awk cal curl date df echo env grep sed sort tail test uname uptime
do
    if ! command -v $c >/dev/null
    then
        echo "WARNING: Command $c not found, check PATH ($PATH)" >&2
    fi
done

# PART III: Check if it is forwarded message or contact #############
#
# If you don't need this functionality, you are free to skip it.
#
# If the bot receives forwarded message or user's contact, it sets
# up three additional variables:
# - BOT_SIDE_TYPE -- type of message: string "user", "contact", "bot" and so on
# - BOT_SIDE_NAME -- name
# - BOT_SIDE_ID -- id, it is useful to set up black/white list in configuration file.
#     if you wish to add some user to black/white list, you can forward his message to bot and
#     figure out the user id. Of cause, you can figure out user id directly from bot's logs.

if test -n "$BOT_SIDE_TYPE"
then
    echo "Message from $BOT_SIDE_TYPE"
    echo "Name: $BOT_SIDE_NAME"
    echo "ID: $BOT_SIDE_ID"
    exit
fi

# PART IV: Commands #################################################
#
# I use magic CMD marker to be able to build help-message automatically

case "CMD_$1" in
    CMD_sup)
        # Simple text messages are sent as is
        #
        echo 'Hi there! 👋'
        ;;
    CMD_date)
        # You can run any command and get output in your Telegram client
        date
        ;;
    CMD_uname)
        # One more example of command
        uname -a
        ;;
    CMD_uptime)
        # And one more
        uptime
        ;;
    CMD_noout)
        # Empty output turns to italic string "empty"
        #
        ;;
    CMD_nothing)
        # Single dot is marker of silence. The bot will reply nothing.
        #
        echo '.'
        ;;
    CMD_args)
        # Try to say to bot
        # - args
        # - args Hello world!
        # and so on
        # You will see how cnbot prepares command arguments:
        # - it converts your message to lower case
        # - save only letters, digits and symbols dot (.), minus (-) and underscore (_)
        # - and consider all the rest symbols as separators
        #
        echo 'Passed args:'
        for a in "$@"
        do
            echo "○ $a"
        done
        echo 'Try:'
        echo 'args'
        echo 'args 1 2 3'
        echo 'args Hello world!'
        ;;
    CMD_env)
        # You can preformat your message, using '%!PRE' control line at the beginning of response.
        #
        # Try this command and check out all available environment variables:
        # BOT_CHAT -- integer chat id
        # BOT_FROM -- integer sender id
        # BOT_FROM_FIRSTNAME -- sender name
        # BOT_NAME -- bot name, according configuration file
        # BOT_SERVER -- server for asynchronous communication (if set up)
        # BOT_TEXT -- original message
        #
        echo '%!PRE'
        env | sort
        ;;
    CMD_cal)
        # One more %!PRE example
        echo '%!PRE'
        cal -h
        ;;
    CMD_calc)
        # The example of how to get access to raw message (BOT_TEXT).
        #
        # ⚠️  WARNING  Be careful, it is very easy to make vulnerability.
        #             Do not use this code in production.
        #
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
            #echo "$prog" # uncomment for debugging
            awk "$prog" 2>&1
        fi
        ;;
    CMD_gologo)
        # To reply by image you have just put image to stdout as is.
        #
        curl -qfs https://golang.org/lib/godoc/images/footer-gopher.jpg
        # curl -qfs https://www.telegram.org/img/t_logo.png # try this if footer-gopher.jpg disappear
        ;;
    CMD_du)
        # You can generate image using APIs, or utilities like RRDtools.
        #
        d="$(df -P -m / | tail -1 | awk '{gsub("[^0-9]", "", $5); print $5","(100-$5)}')"
        # This old fashioned API is deprecated in 2012, however, it is still working
        # https://developers.google.com/chart/image/docs/making_charts
        u="https://chart.googleapis.com/chart?cht=p&chd=t:$d&chs=300x200&chl=Available|Used&chtt=Disk%20usage"
        curl -qfs "$u"
        ;;
    CMD_async)
        # You are free to speak asynchronously using HTTP interface.
        #
        # Here we send three messages in reply to one incoming message.
        #
        # Here you can see
        # - %!MARKDOWN control line
        # - two ways to send asynchronous messages: (1) raw body and (2) multipart/form-data encoding
        #
        # %!MARKDOWN is similar to %!PRE, however
        # - It allows you to use all markdown abbilities
        # - It doesn't escapes control chars for you. So you must to put '\' before every control character by your self.
        #
        # You can also see we specify target user explicitly. So you can send message to other user, chat or even bot, not only in reply to sender.
        #
        mark='%!MARKDOWN'$'\n' # it may surprise you, how we get newline in sh
        # the first way to send async message: multipart/form-data
        curl -qfsX POST -o /dev/null -F to=$BOT_FROM -F msg="${mark}_I'll send you_ *random* _image\.\.\._" $BOT_SERVER
        curl -qfsL https://source.unsplash.com/random/600x400 | curl -qfsX POST -o /dev/null -F to=$BOT_FROM -F msg=@- $BOT_SERVER
        # the second way to send async message: raw data + user_id at the end of url
        echo "${mark}_Are you happy now?_" | curl -qfsX POST -o /dev/null --data-binary @- "http://$BOT_SERVER/$BOT_FROM"
        echo '.' # suppress output processing
        ;;
    CMD_btn)
        # %!CALLBACK is a control line to create inline keyboards
        #
        # There are three ways to use it
        # - %!CALLBACK cmd_name run command
        #   Creates a key with label "run command". That button will start the script with command "cmd_name".
        # - %!CALLBACK one_word
        #   Is literally equal to
        #   %!CALLBACK one_word one_word
        #   label and command are equal
        # - %!CALLBACK
        #   Without arguments the command line starts new line of keys in inline keyboard layout.
        #
        # The following control lines describe three-line keyboard
        # [   ][   ]
        # [        ]
        # [   ][   ]
        #
        echo '%!CALLBACK env do env'
        echo '%!CALLBACK date do date'
        echo '%!CALLBACK'
        echo '%!CALLBACK uname do uname'
        echo '%!CALLBACK'
        echo '%!CALLBACK async get random image'
        echo '%!CALLBACK update update it!'
        echo 'Try buttons'
        ;;
    CMD_update)
        # %!UPDATE control can be used in conjunction with %!CALLBACK and in synchronous messages only
        # it means that the message have to replace the inline-keyboard-message. It allows you to create
        # mutable menus.
        #
        # Here we send the same inline keyboard, but new message text every time.
        #
        echo '%!UPDATE'
        echo '%!CALLBACK env do env'
        echo '%!CALLBACK date do date'
        echo '%!CALLBACK'
        echo '%!CALLBACK uname do uname'
        echo '%!CALLBACK'
        echo '%!CALLBACK async get random image'
        echo '%!CALLBACK update update it!'
        # by the way, you can use JSON APIs with jq assistance
        # curl -qfsL 'http://api.forismatic.com/api/1.0/?method=getQuote&format=json&lang=en' | jq -r '"\(.quoteText) -- \(.quoteAuthor)"'
        curl -qfsL 'http://api.forismatic.com/api/1.0/?method=getQuote&format=text&lang=en'
        ;;
    CMD_help)
        # One more %!MARKDOWN example
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
        # And one more
        cmd="$(echo $1 | sed 's/[-_.]/\\&/g')" # we have to care about [-_.] only because other must-escaped-chars are disallowed by bot
        echo '%!MARKDOWN'
        echo 'This is `cnbot` speaking\.' "I didn't recognize your command '*$cmd*' Try to say '*help*' to me"
        ;;
esac
