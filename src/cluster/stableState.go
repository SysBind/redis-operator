package cluster

type stableState struct {
	cluster *Cluster
}

func (state stableState) boot() error {
	return nil
}

func (state stableState) create() error {
	return nil
}

func (state stableState) scale(count int) error {
	return nil
}

func (state stableState) destroy() error {
	return nil
}
