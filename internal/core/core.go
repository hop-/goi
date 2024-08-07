package core

import (
	"encoding/binary"
	"fmt"
)

func LoadData() error {
	err := loadConsumerGroups()
	if err != nil {
		return err
	}

	return loadTopics()
}

func NewConsumer(consumerName string, groupName string) (*Consumer, error) {
	cg := findConsumerGroupByName(groupName)
	if cg == nil {
		var err error
		cg, err = newConsumerGroup(groupName)
		if err != nil {
			return nil, err
		}
	}

	return newConsumer(consumerName, cg)
}

func RemoveConsumer(name string) error {
	c := findConsumerByName(name)
	if c == nil {
		return fmt.Errorf("unknown consumer %s", name)
	}

	return removeConsumer(c)
}

func NewProducer(name string) (*Producer, error) {
	p := findProducerByName(name)
	if p != nil {
		return nil, fmt.Errorf("producer %s is already exist", name)
	}

	return newProducer(name)
}

func RemoveProducer(name string) error {
	p := findProducerByName(name)
	if p == nil {
		return fmt.Errorf("unknown producer %s", name)
	}

	return removeProducer(p)
}

func NewMessageFromBuff(buff []byte) *Message {
	topicSize := binary.LittleEndian.Uint32(buff[:4])
	topic := string(buff[4 : 4+topicSize])
	content := buff[4+topicSize:]

	return &Message{Topic: topic, Content: content}
}

func NewTopic(name string) (*Topic, error) {
	t := findTopicByName(name)
	if t != nil {
		return nil, fmt.Errorf("topic %s is already exist", name)
	}

	return newTopic(name)
}
