package infra

import (
	"sync"

	"github.com/hop-/goi/internal/core"
	"github.com/hop-/goi/internal/storages"
)

var (
	consumerGroups = make(map[string]*core.ConsumerGroup)
	cgMu           = sync.Mutex{}
)

func newConsumerGroup(name string) (*core.ConsumerGroup, error) {
	cg := core.NewConsumerGroup(name)

	err := addConsumerGroup(cg)

	return cg, err
}

func addConsumerGroup(cg *core.ConsumerGroup) error {
	cgMu.Lock()
	defer cgMu.Unlock()

	err := storages.GetStorage().NewConsumerGroup(cg)
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

	cgs, err := storages.GetStorage().ConsumerGroups()
	if err != nil {
		return err
	}

	for _, cg := range cgs {
		consumerGroups[cg.Name] = &cg
	}

	return nil
}

func findConsumerGroupByName(name string) *core.ConsumerGroup {
	cgMu.Lock()
	defer cgMu.Unlock()

	if cg, ok := consumerGroups[name]; ok {
		return cg
	}

	return nil
}
