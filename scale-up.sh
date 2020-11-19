#!/bin/bash

source config
source common.sh

masters=$(./bin/redis-cli --cluster check 127.0.0.1:6001 | grep "^M: " | wc -l)
count=$((masters*(REPLICAS+1)))

# Master
port=$((6000+count+1))
configure_redis $port
run_redis $port
./bin/redis-cli --cluster add-node 127.0.0.1:$port 127.0.0.1:6001
while ! grep "Cluster state changed: ok" logs/redis-$port.log; do
    sleep 1s
done
master_id=`./bin/redis-cli --cluster check 127.0.0.1:6001 | grep "^M: " | grep $port | cut -d' ' -f 2`
echo "New Master $master_id"

# Replicas
for ((i=$((count+2)); i<=$((count+1+REPLICAS)); i++));
do
  port=$((6000+i))
  configure_redis $port
 run_redis $port
 ./bin/redis-cli --cluster add-node 127.0.0.1:$port 127.0.0.1:6001  --cluster-slave --cluster-master-id $master_id
 while ! grep "MASTER <-> REPLICA sync: Finished with success" logs/redis-$port.log; do
    sleep 1s
 done
done

# Rebalance
 ./bin/redis-cli --cluster fix 127.0.0.1:6001
 sleep 3s
./bin/redis-cli --cluster rebalance 127.0.0.1:6001 --cluster-use-empty-masters
