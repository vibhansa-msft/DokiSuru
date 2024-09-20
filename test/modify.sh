#!/bin/bash

# File to be modified
file="./testdata.1g"
opt=$1
size=$2

if [[ $opt == "trunc" ]]; then
    truncate -s $size $file
    exit 0
else
    # Seek to the offset and modify the data
    for i in {1..100}; do
        offset=$((1 + RANDOM % 2000))
        echo $offset
        dd if=/dev/urandom of="$file" seek="$offset" bs=1M count=1 conv=notrunc
    done
fi