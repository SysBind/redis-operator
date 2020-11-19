#!/bin/bash

source config
source common.sh

masters=$(./bin/redis-cli --cluster check 127.0.0.1:6001 | grep "^M: " | wc -l)
count=$((masters*(REPLICAS+1)))
echo "$masters"
echo "$count"

# ./bin/redis-cli --cluster del-node 127.0.0.01:6001