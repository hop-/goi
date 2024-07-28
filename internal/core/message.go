package core

import "time"

type Message struct {
	Id         *int
	OccurredAt time.Time
	Content    []byte
	Topic      string
}

func newMessage(content []byte, topic *Topic) (*Message, error) {
	m := &Message{
		Content: content,
		Topic:   topic.Name,
	}

	err := addMessage(m)
	return m, err
}

func addMessage(m *Message) error {
	return GetStorage().NewMessage(m)
}
