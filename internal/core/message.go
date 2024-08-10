package core

import (
	"encoding/binary"
)

type Message struct {
	Content []byte
	Topic   string
}

func NewMessage(content []byte, topic *Topic) *Message {
	return &Message{
		Content: content,
		Topic:   topic.Name,
	}
}

func NewMessageFromBuff(buff []byte) *Message {
	topicSize := binary.LittleEndian.Uint32(buff[:4])
	topic := string(buff[4 : 4+topicSize])
	content := buff[4+topicSize:]

	return &Message{Topic: topic, Content: content}
}

func (m *Message) ToBuff() []byte {
	buff := make([]byte, 0, 4+len(m.Topic)+len(m.Content))
	buff = binary.LittleEndian.AppendUint32(buff, uint32(len(m.Topic)))
	buff = append(buff, []byte(m.Topic)...)
	buff = append(buff, m.Content...)

	return buff
}
