package core

import "sync"

type ConsumerGroup struct {
	Id   *int
	Name string
}

var (
	consumerGroups = make(map[string]*ConsumerGroup)
	cgMu           = sync.Mutex{}
)

func newConsumerGroup(name string) (*ConsumerGroup, error) {
	cg := &ConsumerGroup{
		Name: name,
	}

	err := addConsumerGroup(cg)

	return cg, err
}

func addConsumerGroup(cg *ConsumerGroup) error {
	cgMu.Lock()
	defer cgMu.Unlock()

	err := GetStorage().NewConsumerGroup(cg)
	if err != nil {
		return err
	}

	// It should be good to go when storage didn't return an error
	consumerGroups[cg.Name] = cg

	return nil
}

func loadConsumerGroups() error {
	cgMu.Lock()
	defer cgMu.Unlock()

	cgs, err := GetStorage().ConsumerGroups()
	if err != nil {
		return err
	}

	for _, cg := range cgs {
		consumerGroups[cg.Name] = &cg
	}

	return nil
}

func findConsumerGroupByName(name string) *ConsumerGroup {
	cgMu.Lock()
	defer cgMu.Unlock()

	if cg, ok := consumerGroups[name]; ok {
		return cg
	}

	return nil
}
