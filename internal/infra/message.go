package infra

import (
	"github.com/hop-/goi/internal/core"
	"github.com/hop-/goi/internal/storages"
)

func NewMessage(m *core.Message) error {
	// Store message
	err := storages.GetStorage().NewMessage(m)
	if err != nil {
		return err
	}

	// Get all consumer group channels which are waiting for a new message from this topic
	queues := getOnEdgeConsumerGroupChannelsForTopic(m.Topic)
	for _, queue := range queues {
		queue.Push(m)
	}
	// TODO: push message to channels
	return nil
}
