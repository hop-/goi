package core

import (
	"fmt"
	"sync"
)

type Consumer struct {
	Name  string
	group *ConsumerGroup
}

var (
	consumers        = make(map[string]*Consumer)
	consumersByGroup = make(map[string]map[string]*Consumer)
	cMu              = sync.Mutex{}
)

func newConsumer(name string, group *ConsumerGroup) (*Consumer, error) {
	c := &Consumer{
		Name:  name,
		group: group,
	}

	err := addConsumer(c)

	return c, err
}

func addConsumer(c *Consumer) error {
	cMu.Lock()
	defer cMu.Unlock()

	if _, ok := consumers[c.Name]; ok {
		return fmt.Errorf("consumer with %s name is already regisered", c.Name)
	}

	consumers[c.Name] = c
	if cg, ok := consumersByGroup[c.group.Name]; ok {
		if _, ok := cg[c.Name]; ok {
			return fmt.Errorf("consumer with %s name is already registered in cg", c.Name)
		}

		cg[c.Name] = c

		return nil
	}

	consumersByGroup[c.group.Name] = map[string]*Consumer{
		c.Name: c,
	}

	return nil
}

func removeConsumer(c *Consumer) error {
	cMu.Lock()
	defer cMu.Unlock()

	delete(consumers, c.Name)

	if cg, ok := consumersByGroup[c.group.Name]; ok {
		delete(cg, c.Name)
	}

	// no error
	return nil
}

func findConsumerByName(name string) *Consumer {
	cMu.Lock()
	defer cMu.Unlock()

	if c, ok := consumers[name]; ok {
		return c
	}

	return nil
}
