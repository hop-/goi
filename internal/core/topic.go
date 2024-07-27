package core

import "fmt"

type Topic struct {
	Id   *string
	Name string
}

var (
	topics                = make(map[string]Topic)
	topicConsummingGroups = make(map[string]map[string]struct{})
)

func AddTopic(t Topic) error {
	err := GetStorage().NewTopic(t)
	if err != nil {
		return err
	}

	topics[t.Name] = t
	return nil
}

func SubscribeToTopic(topic string, cgName string) error {
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

func UnsubscribeFromTopic(topic string, cgName string) {
	if subs, ok := topicConsummingGroups[topic]; ok {
		delete(subs, cgName)
	}
}

func LoadTopics() error {
	tpcs, err := GetStorage().Topics()
	if err != nil {
		return err
	}

	for _, t := range tpcs {
		topics[t.Name] = t
	}

	return nil
}
