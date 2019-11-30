#!/bin/bash

case "$1" in
    date)
        date
        ;;
    pwd)
        pwd
        ;;
    png)
        convert rose: -resize 100x png:-
        ;;
    jpeg)
        convert rose: -resize 100x jpeg:-
        ;;
    gif)
        convert rose: -resize 100x gif:-
        ;;
    pre)
        echo '```'
        echo 'Art by Donovan Bake'
        echo ' ^ ^'
        echo '(O,O)'
        echo '(   )'
        echo '-"-"---dwb-'
        echo '```'
        ;;
    env)
        echo '```'
        env | sort
        echo '```'
        ;;
    sdate)
        date | curl -s -d '@-' "http://$T_BIND/$T_BOT?to=$T_USER"
        ;;
    spng)
        convert rose: -resize 100x png:- | curl -s --data-binary '@-' "http://$T_BIND/$T_BOT?to=$T_USER"
        ;;
    fpng)
        convert rose: -resize 100x png:- |
            curl -s -F 'photo=@-' -F "text=$(date)" "http://$T_BIND/$T_BOT?to=$T_USER"
        ;;
    fpng2)
        t='*D*ate: `'"$(date)"'`'
        convert rose: -resize 100x png:- |
            curl -s -F 'photo=@-' -F "text=$t" -F "markdown=true" "http://$T_BIND/$T_BOT?to=$T_USER"
        ;;
    md)
        date '+*%Y-%m-%d* `%H:%M:%S`' |
            curl -s -F "text=@-" -F "markdown=true" "http://$T_BIND/$T_BOT?to=$T_USER"
        ;;
    sleep)
        sleep 10
        ;;
    help)
        (
        echo 'Commands'
        echo '```'
        sed -n 's/\([a-z]\))$/\1/p' "$0"
        echo '```'
        ) |
            curl -s -F "text=@-" -F "markdown=true" "http://$T_BIND/$T_BOT?to=$T_USER"
        ;;
    *)
        echo 'Known command "'"$1"'". Say "help" for help'
        ;;
esac
