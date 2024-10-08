package infra

import (
	"sync"

	"github.com/hop-/goi/internal/core"
)

type MessageQueue chan *core.Message

const (
	messageQueueCapacity = 20
)

type ConsumerGroupMessageQueue struct {
	mu    *sync.Mutex
	queue MessageQueue
}

func newConsumerGroupMessageQueue() *ConsumerGroupMessageQueue {
	return &ConsumerGroupMessageQueue{
		queue: make(MessageQueue, messageQueueCapacity),
		mu:    &sync.Mutex{},
	}
}

func (q *ConsumerGroupMessageQueue) IsEmpty() bool {
	return len(q.queue) == 0
}

func (q *ConsumerGroupMessageQueue) IsFull() bool {
	return len(q.queue) == cap(q.queue)
}

func (q *ConsumerGroupMessageQueue) IsHalfEmpty() bool {
	return len(q.queue) <= cap(q.queue)/2
}

func (q *ConsumerGroupMessageQueue) Push(m *core.Message) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsFull() {
		return false
	}

	q.queue <- m

	return true
}

func (q *ConsumerGroupMessageQueue) PushBlocked(m *core.Message) {
	q.queue <- m
}

func (q *ConsumerGroupMessageQueue) PopBlocked() *core.Message {
	return <-q.queue
}
