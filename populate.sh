#!/bin/bash

KEYS_COUNT=1000
source config

for ((i=1; i<=$KEYS_COUNT; i++));
do
  echo "SET key$i VALUE-$i" | ./bin/redis-cli -c -p 6001
done
