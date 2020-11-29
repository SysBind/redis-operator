# Redis Cluster Operator

Status: Alpha

Defines a CRD 'Redis' supporting 'masters', 'replicas' fields and operates a redis cluster according to this spec.

## Building the operator
```make -C src```

## Development

We are using the [kubebuilder](https://book.kubebuilder.io/) framework to scaffold new APIs/CRDs

Since we are assigning same ports for all pods, we'll need a real multi-node environment:
- ```minikube start --nodes 3```

Or (with [kind](https://kind.sigs.k8s.io/))
- ```kind create cluster --config kind-config```

Compile the operator & install CRD to the cluster, than start the operator locally: 
```make -C src && make -C src install && ./src/bin/manager```

Apply a sample Redis spec:
```kubectl apply -f src/config/samples/redis_v1_redis.yaml```

The PODs will be automatically assigned to different nodes because of this port assignment.
TODO: Add option to specify base port so that different deployments of redis will not conflict,
      Or otherwise handle the situation, eg: keep a list of ports in-use.


### Testing Redis Cluster Locally
- `./start.sh`
- `./populate.sh`
- `./scale-up.sh`
- `./scale-down.sh`
- `./stop.sh`

This just tests basic redis-cluster functionality locally, 
Especially fail-over, scale-up with re-balancing & scale-down with re-sharding.
Those files should probably move to a sub-directory.
