package goi

import (
	"fmt"
	"sync"
	"time"

	"github.com/hop-/goi/internal/core"
	"github.com/hop-/goi/internal/network"
	"github.com/hop-/golog"
)

type ConsumerConfig struct {
	Name        *string
	GroupName   *string
	Topic       *string
	Host        *string
	Port        *int
	TcpPort     *int
	TcpFallback *bool
}

type Consumer struct {
	name      string
	groupName string
	topic     string
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
	err = readConfirmation(c)
	if err != nil {
		return err
	}

	// Send producer details
	c.WriteMessage([]byte(name))
	c.WriteMessage([]byte(groupName))

	// Read confirmation
	return readConfirmation(c)
}

func sendCompressionInfo(c *network.Connection) error {
	const compressorType = "none" // TODO: use real compressor
	err := c.WriteMessage([]byte(compressorType))
	if err != nil {
		return err
	}

	// Read confirmation
	return readConfirmation(c)
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

	if config.Topic == nil {
		return nil, fmt.Errorf("topic must be specified")
	}

	fillConsumerDefaults(&config)

	return &Consumer{
		name:      name,
		groupName: groupName,
		topic:     *config.Topic,
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

	err = consumerHandshake(c.conn, c.name, c.groupName)
	if err != nil {
		return err
	}

	err = sendCompressionInfo(c.conn)
	if err != nil {
		return err
	}

	// Ping loop
	go func() {
		errorCounter := 0
		// TODO: add graceful exit
		for {
			time.Sleep(10 * time.Second) // hardcoded
			err := c.conn.Ping()
			if err != nil {
				golog.Error("Failed to ping", err.Error())
				errorCounter += 1
				if errorCounter == 3 {
					break
				}
				continue
			}
			errorCounter = 0
		}
	}()

	return nil
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

	// This is a tmp stupid workaround for quic-go issue
	// https://github.com/quic-go/quic-go/issues/3291
	time.Sleep(3 * time.Second)

	return nil
}

func (c *Consumer) ReadMessage() (string, []byte, error) {
	err := c.conn.WriteSpecialMessage([]byte(network.MessageRequest))
	if err != nil {
		return "", nil, err
	}

	t, buff, err := c.conn.ReadMessage()
	if err != nil {
		return "", nil, err
	}

	if t != network.GeneralMessage {
		return "", nil, fmt.Errorf("unexpected message type")
	}

	// TODO: decompress message
	m := core.NewMessageFromBuff(buff)

	return m.Topic, m.Content, nil
}

// TODO
