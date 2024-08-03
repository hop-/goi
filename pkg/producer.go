package goi

import (
	"sync"

	"github.com/hop-/goi/internal/network"
)

type ProducerConfig struct {
	Host        *string
	Port        *int
	TcpPort     *int
	TcpFallback *bool
}

type Producer struct {
	conn   *network.Connection
	config ProducerConfig
	mu     *sync.Mutex
}

func NewProducer(config ProducerConfig) *Producer {
	return &Producer{
		config: config,
		mu:     &sync.Mutex{},
	}
}

func (p *Producer) Connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	conn, err := connect(*p.config.Host, *p.config.Port, *p.config.TcpPort, *p.config.TcpFallback)
	if err != nil {
		return err
	}

	p.conn = network.NewConnection(conn)

	return nil
}

func (p *Producer) Disconnect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.conn.Close()
	return nil
}

// TODO
