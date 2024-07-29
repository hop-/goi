package handlers

import (
	"github.com/hop-/goi/internal/core"
	"github.com/hop-/goi/internal/network"
)

func handleConsumerHandshake(c *network.Connection) (*core.Consumer, error) {
	// Send confirmation
	err := c.WriteAll(network.OkRes)
	if err != nil {
		return nil, err
	}

	// Read consumer details
	_, consumerName, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}
	_, consumerGroupName, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}
	consumer, err := core.NewConsumer(string(consumerName), string(consumerGroupName))
	if err != nil {
		return nil, err
	}

	// Send confirmation
	err = c.WriteAll(network.OkRes)
	if err != nil {
		return consumer, err
	}

	return consumer, nil
}

func consumerHandler(c *network.Connection) error {
	consumer, err := handleConsumerHandshake(c)
	if consumer != nil {
		// Remove consumer, no matter the handshake status
		defer core.RemoveConsumer(consumer.Name)
	}
	if err != nil {
		return err
	}

	// Consumer main loop
consumerMainLoop:
	for {
		messageType, b, err := c.ReadMessage()
		if err != nil {
			return err
		}

		switch messageType {
		// Handle special codes
		case network.SpecialCode:
			code := b[0]
			// Exit code
			if code == network.ExitCode {
				break consumerMainLoop
			}
		// Handle special messages
		case network.SpecialMessage:
			// Ping pong health check
			if len(b) == 0 {
				c.WriteMessage([]byte{})
			}
		// Handle other
		default:
			// TODO: handle consumer
		}
	}

	return nil
}
