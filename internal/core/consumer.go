package core

type Consumer struct {
	Name  string
	Group *ConsumerGroup
}

func NewConsumer(name string, group *ConsumerGroup) *Consumer {
	return &Consumer{
		Name:  name,
		Group: group,
	}
}
