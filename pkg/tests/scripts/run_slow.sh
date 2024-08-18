#!/bin/sh

trap 'echo trap SIGINT' SIGINT
trap 'echo trap SIGTERM' SIGTERM
trap 'echo trap SIGHUP' SIGHUP
trap 'echo trap SIGQUIT' SIGQUIT
trap 'echo trap EXIT' EXIT
trap 'echo trap ERR' ERR

echo 'start'
sleep 1
echo 'end'
