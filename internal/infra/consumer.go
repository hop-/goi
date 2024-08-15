package infra

import (
	"fmt"
	"sync"

	"github.com/hop-/goi/internal/core"
	"github.com/hop-/golog"
)

var (
	consumers         = make(map[string]*core.Consumer)
	consumersByGroup  = make(map[string]map[string]*core.Consumer)
	topicsByConsumers = make(map[string]map[string]struct{})
	tbsMu             = &sync.Mutex{}
	cMu               = sync.Mutex{}
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

	golog.Debug("New consumer", c.Name)

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

func addTopicToConsumer(consumerName string, topic string) {
	tbsMu.Lock()
	defer tbsMu.Unlock()

	if ts, ok := topicsByConsumers[consumerName]; ok {
		ts[topic] = struct{}{}
	} else {
		topicsByConsumers[consumerName] = map[string]struct{}{
			topic: {},
		}
	}
}

func removeTopicFromConsumer(consumerName string, topic string) {
	tbsMu.Lock()
	defer tbsMu.Unlock()

	if ts, ok := topicsByConsumers[consumerName]; ok {
		delete(ts, topic)
	}
}

func getTopicsByConsumerName(consumerName string) []string {
	tbsMu.Lock()
	defer tbsMu.Unlock()

	if ts, ok := topicsByConsumers[consumerName]; ok {
		topics := make([]string, 0, len(ts))
		for t := range ts {
			topics = append(topics, t)
		}

		return topics
	}

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

func ReadMessage(c *core.Consumer) (*core.Message, error) {
	topics := getTopicsByConsumerName(c.Name)
	if len(topics) == 0 {
		return nil, fmt.Errorf("unable to find any topic for the consumers")
	}

	topic := topics[0]

	channel := getChannel(topic, c.Group.Name)
	if channel == nil {
		return nil, fmt.Errorf("unknown subscription %s for %s topic", c.Group.Name, topic)
	}

	m := channel.PopBlocked()
	return m, nil
}
