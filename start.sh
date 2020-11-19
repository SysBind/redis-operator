#!/bin/bash

source config
source common.sh

get_redis() {
  [ -f redis-${REDIS_RELEASE}.tar.gz ] || wget https://download.redis.io/releases/redis-${REDIS_RELEASE}.tar.gz
  tar xf redis-${REDIS_RELEASE}.tar.gz
  rm -v redis-${REDIS_RELEASE}.tar.gz
  pushd redis-${REDIS_RELEASE} && make && popd
  [ ! -d bin ] && mkdir bin
  cp -v redis-${REDIS_RELEASE}/src/redis-server bin/
  cp -v redis-${REDIS_RELEASE}/src/redis-cli bin/
  rm -rf redis-${REDIS_RELEASE}
}

[ -f bin/redis-server ] || get_redis

addresses=""
for ((i=1; i<=$(( SIZE*(REPLICAS+1) )); i++));
do
  port=$((6000+i))
  configure_redis $port
  run_redis $port
  addresses=$addresses" 127.0.0.1:$port"
done

./bin/redis-cli --cluster create $addresses --cluster-replicas $REPLICAS
sleep $((SIZE*REPLICAS*2))s
./bin/redis-cli --cluster check 127.0.0.1:6001
./bin/redis-cli --cluster rebalance 127.0.0.1:6001 --cluster-use-empty-masters