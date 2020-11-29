package cluster

type creatingState struct {
	cluster *Cluster
}

func (state creatingState) boot() error {
	return nil
}

func (state creatingState) create() error {
	return nil
}

func (state creatingState) scale(count int) error {
	return nil
}

func (state creatingState) destroy() error {
	return nil
}
