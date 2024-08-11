package infra

import (
	"fmt"
	"sync"

	"github.com/hop-/goi/internal/core"
	"github.com/hop-/goi/internal/storages"
)

var (
	topics               = make(map[string]*core.Topic)
	tMu                  = sync.Mutex{}
	subscriptionsByTopic = make(map[string]map[string]*ConsumerGroupMessageQueue)
	cgsMu                = sync.Mutex{}
)

func NewTopic(name string) (*core.Topic, error) {
	t := findTopicByName(name)
	if t != nil {
		return nil, fmt.Errorf("topic %s is already exist", name)
	}

	t = core.NewTopic(name)

	return t, addTopic(t)
}

func addTopic(t *core.Topic) error {
	tMu.Lock()
	defer tMu.Unlock()

	err := storages.GetStorage().NewTopic(t)
	if err != nil {
		return err
	}

	topics[t.Name] = t

	return nil
}

func subscribeToTopic(topic string, cgName string) error {
	cgsMu.Lock()
	defer cgsMu.Unlock()

	if _, ok := topics[topic]; !ok {
		return fmt.Errorf("unknown topic to subscribe %s", topic)
	}

	if subs, ok := subscriptionsByTopic[topic]; ok {
		if _, ok := subs[cgName]; !ok {
			subs[cgName] = newConsumerGroupMessageQueue()
		}
		return nil
	}

	subscriptionsByTopic[topic] = map[string]*ConsumerGroupMessageQueue{
		cgName: newConsumerGroupMessageQueue(),
	}

	return nil
}

func unsubscribeFromTopic(topic string, cgName string) error {
	cgsMu.Lock()
	defer cgsMu.Unlock()

	if subs, ok := subscriptionsByTopic[topic]; ok {
		delete(subs, cgName)
		return nil
	}

	return fmt.Errorf("unknown subscription or topic")
}

func getOnEdgeConsumerGroupChannelsForTopic(topic string) []*ConsumerGroupMessageQueue {
	if subs, ok := subscriptionsByTopic[topic]; ok {
		onEdgeConsumerGroups := make([]*ConsumerGroupMessageQueue, 0, len(subs))
		for _, q := range subs {
			if q.IsOnEdge {
				onEdgeConsumerGroups = append(onEdgeConsumerGroups, q)
			}
		}
		return onEdgeConsumerGroups
	}
	return nil
}

func loadTopics() error {
	tMu.Lock()
	defer tMu.Unlock()

	tpcs, err := storages.GetStorage().Topics()
	if err != nil {
		return err
	}

	for _, t := range tpcs {
		topics[t.Name] = &t
	}

	return nil
}

func findTopicByName(name string) *core.Topic {
	tMu.Lock()
	defer tMu.Unlock()

	if t, ok := topics[name]; ok {
		return t
	}

	return nil
}
