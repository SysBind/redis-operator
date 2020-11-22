#!/bin/bash

source config

masters=$(./bin/redis-cli --cluster check 127.0.0.1:6001 | grep "^M: " | wc -l)
count=$((masters*(REPLICAS+1)))

for ((i=1; i<=$count; i++));
do
  echo ">>> Shutting Down Redis $i.."
  ./bin/redis-cli -p $((6000+i)) "SHUTDOWN"
done

rm -rfv ./conf ./logs
rm -v dump.rdb