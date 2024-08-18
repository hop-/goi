package core

type Message struct {
	Offset  *int64
	Content []byte
	Topic   *Topic
}

func NewMessage(content []byte, topic *Topic) *Message {
	return &Message{
		Offset:  nil,
		Content: content,
		Topic:   topic,
	}
}
