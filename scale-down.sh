#!/bin/bash

source config
source common.sh

get_master_id() {
  idx=$(($1+1))
  ./bin/redis-cli --cluster check 127.0.0.1:6001 | grep "^M: " | head -n$idx | tail -n1  | cut -f2 -d':' | cut -f2 -d' '
}

evacuate_slots() {
  node=$1
  slots=$2
  masters=$(./bin/redis-cli --cluster check 127.0.0.1:6001 | grep "^M: " | wc -l)
  masters=$((masters-1))
  echo "evacuating $slots from $node to remaining $masters masters"

  start_slot=$(echo "$slots" | cut -f1 -d'-')
  end_slot=$(echo "$slots" | cut -f2 -d'-')

  bulksize=$(( (end_slot - start_slot) / masters + 1 ))

  for ((m=0; m<masters; m++)); do
    ./bin/redis-cli --cluster check 127.0.0.1:6001 > logs/check
    if ! grep "All nodes agree about slots configuration." logs/check; then
      echo "Waiting for caluster to stablize.."
      until grep "All nodes agree about slots configuration." logs/check; do
        echo -n ".."
        #./bin/redis-cli --cluster fix 127.0.0.1:6001 > /dev/null || error "could nod fix cluster"
        sleep 3s
        ./bin/redis-cli --cluster check 127.0.0.1:6001 > logs/check
      done
    fi
    id=$(get_master_id $m)
    echo -n "Moving slots to $id.."
    ./bin/redis-cli --cluster reshard  127.0.0.1:6001 --cluster-from $node --cluster-to $id --cluster-slots $bulksize --cluster-yes > /dev/null || error "could not reshard --cluster-from $node --cluster-to $id --cluster-slots $bulksize"
    echo "OK: Moved slots to $id"
  done
}

# Select last master
node=$(./bin/redis-cli --cluster check 127.0.0.1:6001 | grep "^M: " | tail -n1 | cut -f2 -d":" | cut -f2 -d' ')
slots=$(./bin/redis-cli --cluster check 127.0.0.1:6001 | grep -a1 $node | tail -n1 | cut -d":" -f2 | cut -f1 -d' ')


for i in $(echo $slots | sed 's/,/ /g'); do
  i=$(echo $i | tr --delete "[]")
  evacuate_slots $node $i
done

# Remove from cluster (Slaves and than Master)
#slave=$(./bin/redis-cli --cluster check 127.0.0.1:6001 | grep "^S: " | tail -n1 | cut -f2 -d":" | cut -f2 -d' ')
#./bin/redis-cli --cluster del-node 127.0.0.01:6001 $slave
./bin/redis-cli --cluster del-node 127.0.0.01:6001 $node




