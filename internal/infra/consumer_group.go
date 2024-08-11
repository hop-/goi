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

func getOrCreateConsumerGroup(name string) (*core.ConsumerGroup, error) {
	cg := findConsumerGroupByName(name)
	if cg != nil {
		return cg, nil
	}

	cg = core.NewConsumerGroup(name)

	return cg, addConsumerGroup(cg)
}

func addConsumerGroup(cg *core.ConsumerGroup) error {
	cgMu.Lock()
	defer cgMu.Unlock()

	err := storages.GetStorage().NewConsumerGroup(cg)
	if err != nil {
		return err
	}

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
