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
    *)
        echo 'Known commands: pwd, date, png, jpeg, gif. However you said "'"$1"'"'
        ;;
esac
