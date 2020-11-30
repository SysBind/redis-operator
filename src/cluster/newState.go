package cluster

type newState struct {
	cluster *Cluster
}

func (state newState) String() string {
	return "NEW"
}

func (state newState) boot() error {
	log := state.cluster.logger
	// New Redis - Create Statefulset
	if newsts, err := constructStatefulSetForRedis(&state.cluster.spec, state.cluster.scheme); err != nil {
		log.Error(err, "unable to construct statefulset for redis")
		return err
	} else {
		if err := state.cluster.client.Create(state.cluster.ctx, newsts); err != nil {
			log.Error(err, "unable to create Statefulset for Redis", "statefuleset", newsts)
			return err
		}
		log.Info("Created Statefulset for Redis")

		// New Redis - Create Headless Service
		if newsvc, err := constructHeadlessServiceForRedis(&state.cluster.spec, state.cluster.scheme); err != nil {
			log.Error(err, "unable to construct headless service for redis")
			return err
		} else {
			if err := state.cluster.client.Create(state.cluster.ctx, newsvc); err != nil {
				log.Error(err, "unable to create Headless Service for Redis", "service", newsvc)
				return err
			}
		}
		log.Info("Created Headless Service for Redis")
	}
	state.cluster.setState(state.cluster.booting)
	return nil
}

func (state newState) create() error {
	return nil
}

func (state newState) scale(count int) error {
	return nil
}

func (state newState) destroy() error {
	return nil
}
