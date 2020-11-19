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
  echo "generated redis-$port.conf"
}

# run_redis port
## run single instance
run_redis() {
  port=$1
  [ ! -d logs ] && mkdir logs
  ./bin/redis-server conf/redis-$port.conf > logs/redis-$port.log 2>&1 &
  while ! grep "Ready to accept connections" logs/redis-$port.log; do
    sleep 1s
  done
}