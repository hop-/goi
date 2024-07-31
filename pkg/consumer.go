package goi

import "github.com/hop-/goi/internal/network"

type ConsumerConfig struct {
	Host        *string
	Port        *int
	TcpPort     *int
	TcpFallback *bool
}

type Consumer struct {
	conn *network.Connection
}

func NewConsumer(config ConsumerConfig) *Consumer {
	return &Consumer{}
}

func (c *Consumer) Connect() error {
	// TODO
	return nil
}

func (c *Consumer) Disconnect() error {
	// TODO
	return nil
}

// TODO
