package core

type ConsumerGroupState struct {
	Group                *ConsumerGroup
	Topic                *Topic
	Offset               int64
	UnreadMessageOffsets []int64
}

func NewConsumerGroupState(cg *ConsumerGroup, t *Topic, offset int64, messages []int64) *ConsumerGroupState {
	return &ConsumerGroupState{
		Group:                cg,
		Topic:                t,
		Offset:               offset,
		UnreadMessageOffsets: messages,
	}
}
