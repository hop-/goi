package core

import (
	"encoding/binary"
)

type Message struct {
	Offset  *int64
	Content []byte
	Topic   string
}

func NewMessage(content []byte, topic string) *Message {
	return &Message{
		Offset:  nil,
		Content: content,
		Topic:   topic,
	}
}

func NewMessageFromBuff(buff []byte) *Message {
	offset := int64(binary.LittleEndian.Uint64(buff[:8]))
	topicSize := binary.LittleEndian.Uint32(buff[8:12])
	topic := string(buff[12 : 12+topicSize])
	content := buff[12+topicSize:]

	m := NewMessage(content, topic)
	if offset >= 0 {
		m.Offset = &offset
	}

	return m
}

func (m *Message) ToBuff() []byte {
	buff := make([]byte, 0, 12+len(m.Topic)+len(m.Content))

	var offset int64 = -1
	if m.Offset != nil {
		offset = *m.Offset
	}
	buff = binary.LittleEndian.AppendUint64(buff, uint64(offset))

	buff = binary.LittleEndian.AppendUint32(buff, uint32(len(m.Topic)))
	buff = append(buff, []byte(m.Topic)...)
	buff = append(buff, m.Content...)

	return buff
}
