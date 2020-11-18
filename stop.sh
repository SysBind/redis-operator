#!/bin/bash

source config

for ((i=1; i<=$(( SIZE*(REPLICAS+1) )); i++));
do
  echo ">>> Shutting Down Redis $i.."
  ./bin/redis-cli -p $((6000+i)) "SHUTDOWN"
done