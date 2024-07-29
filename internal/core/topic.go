package core

import (
	"fmt"
	"sync"
)

type Topic struct {
	Id   *string
	Name string
}

var (
	topics                = make(map[string]*Topic)
	topicConsummingGroups = make(map[string]map[string]struct{})
	tMu                   = sync.Mutex{}
)

func newTopic(name string) (*Topic, error) {
	t := &Topic{
		Name: name,
	}

	err := addTopic(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func addTopic(t *Topic) error {
	tMu.Lock()
	defer tMu.Unlock()

	err := GetStorage().NewTopic(t)
	if err != nil {
		return err
	}

	topics[t.Name] = t
	return nil
}

func subscribeToTopic(topic string, cgName string) error {
	tMu.Lock()
	defer tMu.Unlock()

	if _, ok := topics[topic]; !ok {
		return fmt.Errorf("unknown topic to subscribe %s", topic)
	}

	if subs, ok := topicConsummingGroups[topic]; ok {
		subs[cgName] = struct{}{}
		return nil
	}

	topicConsummingGroups[topic] = map[string]struct{}{
		cgName: {},
	}

	return nil
}

func unsubscribeFromTopic(topic string, cgName string) {
	tMu.Lock()
	defer tMu.Unlock()

	if subs, ok := topicConsummingGroups[topic]; ok {
		delete(subs, cgName)
	}
}

func loadTopics() error {
	tMu.Lock()
	defer tMu.Unlock()

	tpcs, err := GetStorage().Topics()
	if err != nil {
		return err
	}

	for _, t := range tpcs {
		topics[t.Name] = &t
	}

	return nil
}

func findTopicByName(name string) *Topic {
	tMu.Lock()
	defer tMu.Unlock()

	if t, ok := topics[name]; ok {
		return t
	}

	return nil
}
