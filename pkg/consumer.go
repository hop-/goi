package goi

import (
	"fmt"
	"sync"
	"time"

	"github.com/hop-/goi/internal/network"
)

type ConsumerConfig struct {
	Name        *string
	GroupName   *string
	Host        *string
	Port        *int
	TcpPort     *int
	TcpFallback *bool
}

type Consumer struct {
	name      string
	groupName string
	conn      *network.Connection
	config    ConsumerConfig
	mu        *sync.Mutex
}

func fillConsumerDefaults(c *ConsumerConfig) {
	if c.Host == nil {
		c.Host = &defaultHost
	}
	if c.Port == nil {
		c.Port = &defaultPort
	}
	if c.TcpPort == nil {
		c.TcpPort = c.Port
	}
	if c.TcpFallback == nil {
		c.TcpFallback = &defaultFallback
	}
}

func consumerHandshake(c *network.Connection, name string, groupName string) error {
	// Send client type
	err := c.WriteAll(network.ConsumerTypeMessage)
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
	c.WriteMessage([]byte(groupName))

	// Read confirmation
	err = c.ReadAll(smallBuff)
	if err != nil {
		return err
	} else if smallBuff[0] != network.OkResCode {
		return rejectErr
	}

	return nil
}

func NewConsumer(config ConsumerConfig) (*Consumer, error) {
	var name, groupName string
	var err error
	if config.Name == nil {
		name, err = randomUuidAsString()
		if err != nil {
			return nil, err
		}
	} else {
		name = *config.Name
	}
	if config.GroupName == nil {
		groupName, err = randomUuidAsString()
		if err != nil {
			return nil, err
		}
	} else {
		groupName = *config.GroupName
	}

	fillConsumerDefaults(&config)

	return &Consumer{
		name:      name,
		groupName: groupName,
		config:    config,
		mu:        &sync.Mutex{},
	}, nil
}

func (c *Consumer) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, err := connect(*c.config.Host, *c.config.Port, *c.config.TcpPort, *c.config.TcpFallback)
	if err != nil {
		return err
	}

	c.conn = network.NewConnection(conn)

	return consumerHandshake(c.conn, c.name, c.groupName)
}

func (c *Consumer) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	defer c.conn.Close()

	// Send exit code
	err := c.conn.WriteSpecialCode(network.ExitCode)
	if err != nil {
		return err
	}

	// This is tmp stupid workaround for quic-go issue
	// https://github.com/quic-go/quic-go/issues/3291
	time.Sleep(3 * time.Second)

	return nil
}

// TODO
