#!/bin/bash

echo '%!PRE'

# (12+1)*315 = 4095 chars
for i in {1..315}
do
    echo "⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘" # 12+1 chars (including "\n")
done
echo 12 # exactly 4097 chars after trimming spaces
