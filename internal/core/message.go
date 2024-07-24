package core

import "time"

var (
	// Message statuses
	SentStatus     = "sent"
	AcceptedStatus = "accepted"
	PendingStatus  = "pending"
)

type Message struct {
	Id         *int
	OccurredAt time.Time
	Content    []byte
	TOPIC      string
}
