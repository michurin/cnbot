#!/bin/bash

# This script is for processing "tasks" those are scheduled by script.sh
#
# It is assumed, that this script will be run periodically.
# For example, by crontab line:
# * * * * * cron.sh

BOTS_ASYNC_INTERFACE='http://127.0.0.1:9092'

bound="$(( $(date +%s) - 100 ))"

cd "$(dirname "$0")/data"

for d in $(ls -1)
do
    read t u a b c <<< "$(sed 's/-/ /g' <<< "$d")"
    if test "$t" -lt "$bound"
    then
        (
            echo '%!MARKDOWN'
            echo "I have done the task *$a \\* $b \\* $c* and found out that"
            echo "*$a \\* $b \\* $c \\= __$(($a*$b*$c))__*"
        ) |
        curl -qfsX POST -o /dev/null -F to=$u -F msg=@- "$BOTS_ASYNC_INTERFACE"
        rm "$d"
    fi
done
