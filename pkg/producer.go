package goi

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/hop-/goi/internal/network"
)

type ProducerConfig struct {
	Name        *string
	Host        *string
	Port        *int
	TcpPort     *int
	TcpFallback *bool
}

type Producer struct {
	name   string
	conn   *network.Connection
	config ProducerConfig
	mu     *sync.Mutex
}

func randomUuidAsString() (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return uid.String(), nil
}

func producerHandshake(c *network.Connection, name string) error {
	// Send client type
	err := c.WriteAll(network.ProducerTypeMessage)
	if err != nil {
		return err
	}

	// Read confirmation
	rejectErr := fmt.Errorf("handshake rejected from server")
	smallBuff := make([]byte, 1)
	err = c.ReadAll(smallBuff)
	if err != nil {
		return err
	} else if smallBuff[0] != network.OkResCode {
		return rejectErr
	}

	// Send producer details
	c.WriteMessage([]byte(name))

	// Read confirmation
	err = c.ReadAll(smallBuff)
	if err != nil {
		return err
	} else if smallBuff[0] != network.OkResCode {
		return rejectErr
	}

	return nil
}

func NewProducer(config ProducerConfig) (*Producer, error) {
	var name string
	var err error
	if config.Name != nil {
		name, err = randomUuidAsString()
		if err != nil {
			return nil, err
		}
	} else {
		name = *config.Name
	}

	return &Producer{
		name:   name,
		config: config,
		mu:     &sync.Mutex{},
	}, nil
}

func (p *Producer) Connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	conn, err := connect(*p.config.Host, *p.config.Port, *p.config.TcpPort, *p.config.TcpFallback)
	if err != nil {
		return err
	}

	p.conn = network.NewConnection(conn)

	return producerHandshake(p.conn, p.name)
}

func (p *Producer) Disconnect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Send exit code
	p.conn.WriteSpecialCode(network.ExitCode)

	p.conn.Close()
	return nil
}

// TODO
