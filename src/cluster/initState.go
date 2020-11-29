package cluster

type initState struct {
	cluster *Cluster
}

func (state initState) boot() error {
	return nil
}

func (state initState) create() error {
	return nil
}

func (state initState) scale(count int) error {
	return nil
}

func (state initState) destroy() error {
	return nil
}
