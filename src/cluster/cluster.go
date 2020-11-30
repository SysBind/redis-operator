package cluster

import (
	"context"
	"github.com/go-logr/logr"
	redisv1 "gitlab.sysbind.biz/operators/redis-operator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// using the 'state' design pattern to easily handle cluster logic
// (Modeled as states and transitions)
// see https://golangbyexample.com/state-design-pattern-go/
type state interface {
	boot() error
	create() error
	scale(count int) error
	destroy() error
}

type Cluster struct {
	new          state // still no objects created for this spec
	booting      state // still loading statefuls set pods and other object
	init         state // have pods and other objects, need to initiate the cluster (--cluster create)
	creating     state // We have issued --cluster create, waiting for it to stabilize
	stable       state // cluster is up and stable
	destroying   state // cluster is being destroyed
	currentState state

	spec   redisv1.Redis
	client client.Client
	logger logr.Logger
	scheme *runtime.Scheme
	ctx    context.Context
}

func NewCluster(spec redisv1.Redis, client client.Client, scheme *runtime.Scheme, ctx context.Context, logger logr.Logger) *Cluster {
	cluster := &Cluster{spec: spec, client: client, scheme: scheme, ctx: ctx, logger: logger}

	// wire-up states
	cluster.new = newState{cluster: cluster}
	cluster.booting = bootingState{cluster: cluster}
	cluster.init = initState{cluster: cluster}
	cluster.creating = creatingState{cluster: cluster}
	cluster.stable = stableState{cluster: cluster}
	cluster.destroying = destroyingState{cluster: cluster}
	cluster.setState(cluster.new)

	return cluster
}

func (c *Cluster) setState(state state) {
	if c.currentState != nil {
		c.logger.Info("Redis Cluster setState ", "prev", c.currentState)
	} else {
		c.logger.Info("Current state is NIL")
	}
	c.logger.Info("Redis Cluster setState ", "new", state)
	c.currentState = state
}

func (c *Cluster) Boot() error {
	return c.currentState.boot()
}

func (c *Cluster) Create() error {
	return c.currentState.create()
}

func (c *Cluster) Scale(count int) error {
	return c.currentState.scale(count)
}
