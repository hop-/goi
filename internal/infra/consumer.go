package infra

import (
	"fmt"
	"sync"

	"github.com/hop-/goi/internal/core"
)

var (
	consumers        = make(map[string]*core.Consumer)
	consumersByGroup = make(map[string]map[string]*core.Consumer)
	cMu              = sync.Mutex{}
)

func NewConsumer(consumerName string, groupName string) (*core.Consumer, error) {
	cg, err := getOrCreateConsumerGroup(groupName)
	if err != nil {
		return nil, err
	}

	c := core.NewConsumer(consumerName, cg)

	return c, addConsumer(c)
}

func RemoveConsumer(name string) error {
	c := findConsumerByName(name)
	if c == nil {
		return fmt.Errorf("unknown consumer %s", name)
	}

	return removeConsumer(c)
}

func GetIdleConsumersByTopicName(topic string) ([]*core.Consumer, error) {
	cs := make([]*core.Consumer, 0)

	return cs, nil
}

func addConsumer(c *core.Consumer) error {
	cMu.Lock()
	defer cMu.Unlock()

	if _, ok := consumers[c.Name]; ok {
		return fmt.Errorf("consumer with %s name is already regisered", c.Name)
	}

	consumers[c.Name] = c
	if cg, ok := consumersByGroup[c.Group.Name]; ok {
		if _, ok := cg[c.Name]; ok {
			return fmt.Errorf("consumer with %s name is already registered in cg", c.Name)
		}

		cg[c.Name] = c

		return nil
	}

	consumersByGroup[c.Group.Name] = map[string]*core.Consumer{
		c.Name: c,
	}

	return nil
}

func removeConsumer(c *core.Consumer) error {
	cMu.Lock()
	defer cMu.Unlock()

	delete(consumers, c.Name)

	if cg, ok := consumersByGroup[c.Group.Name]; ok {
		delete(cg, c.Name)
	}

	// no error
	return nil
}

func findConsumerByName(name string) *core.Consumer {
	cMu.Lock()
	defer cMu.Unlock()

	if c, ok := consumers[name]; ok {
		return c
	}

	return nil
}
