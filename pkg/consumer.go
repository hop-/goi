package goi

import (
	"sync"

	"github.com/hop-/goi/internal/network"
)

type ConsumerConfig struct {
	Host        *string
	Port        *int
	TcpPort     *int
	TcpFallback *bool
}

type Consumer struct {
	conn   *network.Connection
	config ConsumerConfig
	mu     *sync.Mutex
}

func NewConsumer(config ConsumerConfig) *Consumer {
	return &Consumer{
		config: config,
		mu:     &sync.Mutex{},
	}
}

func (c *Consumer) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, err := connect(*c.config.Host, *c.config.Port, *c.config.TcpPort, *c.config.TcpFallback)
	if err != nil {
		return err
	}

	c.conn = network.NewConnection(conn)

	return nil
}

func (c *Consumer) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.conn.Close()
	return nil
}

// TODO
