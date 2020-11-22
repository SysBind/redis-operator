#!/bin/bash

source config
source common.sh

while ./bin/redis-cli --cluster check 127.0.0.1:6001 | grep "^M: " > /dev/null;
do
  node=$(./bin/redis-cli --cluster check 127.0.0.1:6001 | grep "^M: " | tail -n1 | cut -f2 -d":" | cut -f2 -d' ')
  port=$(./bin/redis-cli --cluster check 127.0.0.1:6001 | grep "^M: " | tail -n1 | cut -f3 -d":")
  echo "remove_replicas for $port"
  remove_replicas $node
  echo "shutdown master: shutdown_redis $port"
  shutdown_redis $port
done


rm -rfv ./conf ./logs
rm -v dump.rdb