package core

type ConsumerGroup struct {
	Name string
}

func NewConsumerGroup(name string) *ConsumerGroup {
	return &ConsumerGroup{
		Name: name,
	}
}
