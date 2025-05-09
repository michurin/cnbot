#!/bin/bash

# (12+1)*315 = 4095 chars
for i in {1..315}
do
    echo "⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘⌘" # 12+1 chars (including "\n")
done
echo 1 # exactly 4096 chars after trimming spaces
