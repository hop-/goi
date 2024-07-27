package core

import "time"

type Message struct {
	Id         *int
	OccurredAt time.Time
	Content    []byte
	Topic      string
}

func AddMessage(m Message) error {
	return GetStorage().NewMessage(m)
}
