package handlers

import (
	"encoding/binary"

	"github.com/hop-/goi/internal/core"
	"github.com/hop-/goi/internal/infra"
	"github.com/hop-/goi/internal/network"
	"github.com/hop-/golog"
)

func handleConsumerHandshake(c *network.Connection) (*core.Consumer, func(), error) {
	// Send confirmation
	err := c.WriteAll(network.OkRes)
	if err != nil {
		return nil, nil, err
	}

	// Read consumer details
	_, consumerName, err := c.ReadMessage()
	if err != nil {
		return nil, nil, err
	}
	_, consumerGroupName, err := c.ReadMessage()
	if err != nil {
		return nil, nil, err
	}
	consumer, err := infra.NewConsumer(string(consumerName), string(consumerGroupName))
	if err != nil {
		return nil, nil, err
	}

	// Read topics to subscribe
	var topicsCount int32
	err = binary.Read(c, binary.LittleEndian, &topicsCount)
	if err != nil {
		return consumer, nil, err
	}

	topics := make([]string, 0, topicsCount)
	for i := 0; i < int(topicsCount); i++ {
		_, topic, err := c.ReadMessage()
		if err != nil {
			return consumer, nil, err
		}

		topics = append(topics, string(topic))
	}

	golog.Debugf("Topics for the consumer %s are %v", consumer.Name, topics)

	// Subscribe
	validTopics := make([]string, 0, topicsCount)
	for _, topic := range topics {
		err = infra.Subscribe(string(topic), consumer)
		if err != nil {
			golog.Error("Failed to subscribe to the topic", err)
			continue
		}
		validTopics = append(validTopics, topic)
	}

	unsubscribeFunc := func() {
		for _, topic := range validTopics {
			err = infra.Unsubscribe(string(topic), consumer)
			if err != nil {
				golog.Error("Failed to unsubscribe from the topic", err)
				continue
			}
		}
	}

	// Send confirmation
	err = c.WriteAll(network.OkRes)
	if err != nil {
		return consumer, unsubscribeFunc, err
	}

	return consumer, unsubscribeFunc, nil
}

func consumerHandler(c *network.Connection) error {
	consumer, ufunc, err := handleConsumerHandshake(c)
	if consumer != nil {
		// Remove consumer, no matter the handshake status
		defer infra.RemoveConsumer(consumer.Name)
	}
	if ufunc != nil {
		defer ufunc()
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
			specialMessage := string(b)
			if specialMessage == network.MessageRequest {
				// Request a message
				_ = compressor
				// TODO: handle message request
			}
		}
	}

	return nil
}
