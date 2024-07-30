package handlers

import (
	"github.com/hop-/goi/internal/core"
	"github.com/hop-/goi/internal/network"
)

func handleProducerHandshake(c *network.Connection) (*core.Producer, error) {
	// Send confirmation
	c.WriteAll(network.OkRes)

	// Read producer details
	_, producerName, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}

	producer, err := core.NewProducer(string(producerName))
	if err != nil {
		return nil, err
	}

	// Send confirmation
	err = c.WriteAll(network.OkRes)
	if err != nil {
		return producer, err
	}

	return producer, nil
}

func producerHandler(c *network.Connection) error {
	producer, err := handleProducerHandshake(c)
	if producer != nil {
		// Remove producer, no matter the handshake status
		defer core.RemoveProducer(producer.Name)
	}
	if err != nil {
		return err
	}

	// Producer main loop
producerMainLoop:
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
				break producerMainLoop
			}
		// Handle special messages
		case network.SpecialMessage:
			// Ping pong health check
			if len(b) == 0 {
				c.WriteMessage([]byte{})
			}
		// Handle other
		default:
			// TODO: handle producer
		}
	}

	return nil
}
