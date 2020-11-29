package cluster

type destroyingState struct {
	cluster *Cluster
}

func (state destroyingState) boot() error {
	return nil
}

func (state destroyingState) create() error {
	return nil
}

func (state destroyingState) scale(count int) error {
	return nil
}

func (state destroyingState) destroy() error {
	return nil
}
