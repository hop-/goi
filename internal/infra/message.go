package infra

import (
	"encoding/binary"
	"fmt"

	"github.com/hop-/goi/internal/core"
	"github.com/hop-/goi/internal/storages"
	"github.com/hop-/golog"
)

func NewMessageFromBuff(buff []byte) (*core.Message, error) {
	offset := int64(binary.LittleEndian.Uint64(buff[:8]))
	topicSize := binary.LittleEndian.Uint32(buff[8:12])
	topicName := string(buff[12 : 12+topicSize])
	topic := findTopicByName(topicName)
	if topic == nil {
		return nil, fmt.Errorf("message on unknown topic %s", topicName)
	}
	content := buff[12+topicSize:]

	m := core.NewMessage(content, topic)
	if offset >= 0 {
		m.Offset = &offset
	}

	return m, nil
}

func MessageToBuff(m *core.Message) ([]byte, error) {
	if m.Topic == nil {
		return nil, fmt.Errorf("topic must be set")
	}
	buff := make([]byte, 0, 12+len(m.Topic.Name)+len(m.Content))

	var offset int64 = -1
	if m.Offset != nil {
		offset = *m.Offset
	}
	buff = binary.LittleEndian.AppendUint64(buff, uint64(offset))

	buff = binary.LittleEndian.AppendUint32(buff, uint32(len(m.Topic.Name)))
	buff = append(buff, []byte(m.Topic.Name)...)
	buff = append(buff, m.Content...)

	return buff, nil
}
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
	queues := getConsumerGroupChannelsForTopic(m.Topic.Name)
	golog.Debug(len(queues), "Queue(s) to push")
	for cgName, queue := range queues {
		pushed := queue.Push(m)

		if !pushed {
			go loadMessagesFromStorage(cgName, m.Topic.Name, queue)
		}
	}

	return nil
}

func loadMessagesFromStorage(consumerGroupName string, topicName string, queue *ConsumerGroupMessageQueue) {
	for {
		message, err := getNextMessageFromStorage(consumerGroupName, topicName)
		if err != nil {
			golog.Error("Failed to retrieve a message from the storage", err.Error())
			// TODO: handle the error properly
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

	// TODO: get consumer group stats and next message offset
	offset := 1

	return storages.GetStorage().MessageByTopicAndOffset(t, int64(offset))
}
