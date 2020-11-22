
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


error() {
  echo $1
  exit 2
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

shutdown_redis() {
  port=$1
  echo ">>> Shutting Down Redis on $port.."
  ./bin/redis-cli -p $port "SHUTDOWN"
  rm -v conf/redis-$port.conf
  rm -v conf/nodes-$port.conf
}

# remove replicas of master
remove_replicas() {
  node=$1
  while ./bin/redis-cli --cluster check 127.0.0.1:6001  | grep -b2 "replicates $node" > /dev/null; do
    replica=$(./bin/redis-cli --cluster check 127.0.0.1:6001  | grep -b2 "replicates $node" | head -n1 | cut -f2 -d':' | cut -f2 -d' ')
    port=$(./bin/redis-cli --cluster check 127.0.0.1:6001  | grep -b2 "replicates $node" | head -n1 | cut -f3 -d':')
    echo "removing replica $replica"
    /bin/redis-cli --cluster del-node 127.0.0.01:6001 $replica
    shutdown_redis $port
  done
}
