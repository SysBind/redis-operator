# Redis Cluster Operator

## Testing Redis Cluster Locally
- `./start.sh`
- `./populate.sh`
- `./scale-up.sh`
- `./scale-down.sh`
- `./stop.sh`

This just tests basic redis-cluster functionality locally, 
Especially fail-over, scale-up with re-balancing & scale-down with re-sharding.
Those file should probably move to a sub-directory.

## Building the operator
```make -C src```

## Development

We are using the [kubebuilder](https://book.kubebuilder.io/) framework to scaffold new APIs/CRDs


