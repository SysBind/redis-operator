#!/bin/bash

source config

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

# configure_redis port
## configure single instance
configure_redis() {
  [ -d conf ] || mkdir conf
  port=$1
  cat <<EOF > conf/redis-$port.conf
  port $port
  cluster-enabled yes
  cluster-config-file conf/nodes-$port.conf
  appendonly no
  save ""
EOF
  echo "generated redis-$num.conf"
}

[ -f bin/redis-server ] || get_redis

addresses=""
for ((i=1; i<=$(( SIZE*(REPLICAS+1) )); i++));
do
  port=$((6000+i))
  configure_redis $port
  ./bin/redis-server conf/redis-$port.conf &
  addresses=$addresses" 127.0.0.1:$port"
done

./bin/redis-cli --cluster create $addresses --cluster-replicas $REPLICAS
./bin/redis-cli --cluster check 127.0.0.1:6001