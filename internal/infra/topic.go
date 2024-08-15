package infra

import (
	"fmt"
	"sync"

	"github.com/hop-/goi/internal/core"
	"github.com/hop-/goi/internal/storages"
	"github.com/hop-/golog"
)

var (
	topics                = make(map[string]*core.Topic)
	tMu                   = sync.Mutex{}
	subscriptionsByTopics = make(map[string]map[string]*ConsumerGroupMessageQueue)
	cgsMu                 = sync.Mutex{}
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

	golog.Debug("New topic", t.Name)

	return nil
}

func subscribeToTopic(topic string, cgName string) error {
	cgsMu.Lock()
	defer cgsMu.Unlock()

	if _, ok := topics[topic]; !ok {
		return fmt.Errorf("unknown topic to subscribe %s", topic)
	}

	if subs, ok := subscriptionsByTopics[topic]; ok {
		if _, ok := subs[cgName]; !ok {
			q := newConsumerGroupMessageQueue()
			loadMessagesFromStorage(cgName, topic, q)
			subs[cgName] = q
		}

		return nil
	}

	q := newConsumerGroupMessageQueue()
	loadMessagesFromStorage(cgName, topic, q)

	subscriptionsByTopics[topic] = map[string]*ConsumerGroupMessageQueue{
		cgName: q,
	}

	return nil
}

func unsubscribeFromTopic(topic string, cgName string) error {
	cgsMu.Lock()
	defer cgsMu.Unlock()

	if subs, ok := subscriptionsByTopics[topic]; ok {
		delete(subs, cgName)
		return nil
	}

	return fmt.Errorf("unknown subscription or topic")
}

func getConsumerGroupChannelsForTopic(topic string) map[string]*ConsumerGroupMessageQueue {
	cgsMu.Lock()
	defer cgsMu.Unlock()

	if subs, ok := subscriptionsByTopics[topic]; ok {
		return subs
	}

	return nil
}

func getChannel(topic string, consumerGroupName string) *ConsumerGroupMessageQueue {
	cgsMu.Lock()
	defer cgsMu.Unlock()

	if subs, ok := subscriptionsByTopics[topic]; ok {
		if channel, ok := subs[consumerGroupName]; ok {
			return channel
		}
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

	// Add test topic if not exist
	if _, ok := topics["test"]; !ok {
		topics["test"] = core.NewTopic("test")
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
