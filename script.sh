#!/bin/sh

if test "x$1" = 'xdot'
then
    echo ' . '
    exit 0
fi

if test "x$1" = 'xempty'
then
    exit 0
fi

if test "x$1" = 'xenv'
then
    echo '```'
    env
    pwd
    ls
    whereis curl
    ps uxwww -p $$
    echo '```'
    exit 0
fi

if test "x$1" = 'xpng'
then
    convert rose: -resize 200x png:-
    exit 0
fi

if test "x$1" = 'xcurl'
then
    url="http://localhost:9090/$BOT_NAME/to/$BOT_FROM"
    echo "I'm sending image to you now..." |
        curl -qsX POST -o /dev/null --data-binary @- "$url"
    convert rose: -resize 200x png:- |
        curl -qsX POST -o /dev/null --data-binary @- "$url"
    echo "enjoy!"
    exit 0
fi

for i in "$@"
do
    echo "= $i"
done
