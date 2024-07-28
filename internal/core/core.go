package core

import "fmt"

func Init() error {
	err := loadConsumerGroups()
	if err != nil {
		return err
	}

	return loadTopics()
}

func NewConsumer(consumerName string, groupName string) error {
	cg := findConsumerGroupByName(groupName)
	if cg == nil {
		var err error
		cg, err = newConsumerGroup(groupName)
		if err != nil {
			return err
		}
	}

	_, err := newConsumer(consumerName, cg)
	return err
}

func RemoveConsumer(name string) error {
	c := findConsumerByName(name)
	if c == nil {
		return fmt.Errorf("unknown consumer %s", name)
	}

	return removeConsumer(c)
}

func NewMessage(content []byte, topicName string) error {
	topic := findTopicByName(topicName)
	if topic == nil {
		return fmt.Errorf("unknown topic %s", topicName)
	}

	_, err := newMessage(content, topic)
	return err
}
