#!/bin/bash

source config
source common.sh

get_redis() {
  [ -f redis-${REDIS_RELEASE}.tar.gz ] || curl -LON https://download.redis.io/releases/redis-${REDIS_RELEASE}.tar.gz
  tar xf redis-${REDIS_RELEASE}.tar.gz
  rm -v redis-${REDIS_RELEASE}.tar.gz
  pushd redis-${REDIS_RELEASE} && make && popd
  [ ! -d bin ] && mkdir bin
  cp -v redis-${REDIS_RELEASE}/src/redis-server bin/
  cp -v redis-${REDIS_RELEASE}/src/redis-cli bin/
  rm -rf redis-${REDIS_RELEASE}
}

[ -f bin/redis-server ] || get_redis

[ -d ./conf ] && error "cluster already started, please ./stop.sh first"

addresses=""
for ((i=1; i<=$(( SIZE*(REPLICAS+1) )); i++));
do
  port=$((6000+i))
  configure_redis $port
  run_redis $port
  addresses=$addresses" 127.0.0.1:$port"
done

./bin/redis-cli --cluster create $addresses --cluster-replicas $REPLICAS --cluster-yes
