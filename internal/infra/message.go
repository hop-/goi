package infra

import (
	"fmt"

	"github.com/hop-/goi/internal/core"
	"github.com/hop-/goi/internal/storages"
	"github.com/hop-/golog"
)

func ProcessNewMessage(m *core.Message) error {
	// Store message
	err := storages.GetStorage().NewMessage(m)
	if err != nil {
		return err
	}

	// TODO: update message offset
	var offset int64 = 1 // hardcoded
	m.Offset = &offset

	// Get all consumer group channels which are waiting for a new message from this topic
	queues := getConsumerGroupChannelsForTopic(m.Topic)
	golog.Debug(len(queues), "Queue(s) to push")
	for cgName, queue := range queues {
		pushed := queue.Push(m)

		if !pushed {
			go loadMessagesFromStorage(cgName, m.Topic, queue)
		}
	}

	return nil
}

func loadMessagesFromStorage(consumerGroupName string, topicName string, queue *ConsumerGroupMessageQueue) {
	for {
		message, err := getNextMessageFromStorage(consumerGroupName, topicName)
		if err != nil {
			golog.Error("Failed to retrieve a message from the storage", err.Error())
			// TODO: manage the error
			continue
		} else if message == nil {
			// No more message to retrieve from the storage
			break
		}

		queue.PushBlocked(message)
	}
}

func getNextMessageFromStorage(consumerGroupName string, topicName string) (*core.Message, error) {
	cg := findConsumerGroupByName(consumerGroupName)
	if cg == nil {
		return nil, fmt.Errorf("unknown consumer group %s", consumerGroupName)
	}

	t := findTopicByName(topicName)
	if t == nil {
		return nil, fmt.Errorf("unknown topic %s", topicName)
	}

	return storages.GetStorage().NextMessageForConsumerGroup(cg, t)
}
