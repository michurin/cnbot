#!/bin/sh

if test "x$1" = 'xpng'
then
    convert rose: -resize 200x png:-
    exit 0
fi

for i in "$@"
do
    echo "= $i"
done