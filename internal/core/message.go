package core

import (
	"encoding/binary"
)

type Message struct {
	Content []byte
	Topic   string
}

func newMessage(content []byte, topic *Topic) (*Message, error) {
	m := &Message{
		Content: content,
		Topic:   topic.Name,
	}

	err := addMessage(m)
	return m, err
}

func (m *Message) ToBuff() []byte {
	buff := make([]byte, 0, 4+len(m.Topic)+len(m.Content))
	buff = binary.LittleEndian.AppendUint32(buff, uint32(len(m.Topic)))
	buff = append(buff, []byte(m.Topic)...)
	buff = append(buff, m.Content...)

	return buff
}

func addMessage(m *Message) error {
	return GetStorage().NewMessage(m)
}
