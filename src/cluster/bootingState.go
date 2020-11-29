package cluster

type bootingState struct {
	cluster *Cluster
}

func (state bootingState) String() string {
	return "BOOTING"
}

func (state bootingState) boot() error {
	return nil
}

func (state bootingState) create() error {
	return nil
}

func (state bootingState) scale(count int) error {
	return nil
}

func (state bootingState) destroy() error {
	return nil
}
