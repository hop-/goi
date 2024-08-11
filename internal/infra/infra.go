package infra

import (
	"github.com/hop-/goi/internal/core"
)

func LoadData() error {
	err := loadConsumerGroups()
	if err != nil {
		return err
	}

	return loadTopics()
}

func Subscribe(topic string, c *core.Consumer) error {
	return subscribeToTopic(topic, c.Group.Name)
}

func Unsubscribe(topic string, c *core.Consumer) error {
	return unsubscribeFromTopic(topic, c.Group.Name)
}
