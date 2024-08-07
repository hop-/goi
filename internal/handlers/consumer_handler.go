package handlers

import (
	"github.com/hop-/goi/internal/core"
	"github.com/hop-/goi/internal/network"
	"github.com/hop-/golog"
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

	compressor, err := getCompressor(c)
	if err != nil {
		return err
	}

	golog.Infof("New consumer %s accepted", consumer.Name)

	// Consumer main loop
consumerMainLoop:
	for {
		golog.Debug("Waiting message from consumer", consumer.Name)
		messageType, b, err := c.ReadMessage()
		if err != nil {
			return err
		}
		golog.Debug("New message has been received from consumer", consumer.Name)

		switch messageType {
		// Handle special codes
		case network.SpecialCode:
			code := b[0]
			// Exit code
			if code == network.ExitCode {
				golog.Info("Received exit code from consumer", consumer.Name)
				break consumerMainLoop
			}
		// Ping pong health check
		case network.PingMessage:
			golog.Debug("Received ping from consumer", consumer.Name)
			continue consumerMainLoop
		// Special messages
		case network.SpecialMessage:
			switch string(b) {
			// requset a message
			case network.MessageRequest:
				_ = compressor
				// TODO: handle message request
			}
		}
	}

	return nil
}
