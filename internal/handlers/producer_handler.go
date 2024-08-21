package handlers

import (
	"github.com/hop-/goi/internal/compressors"
	"github.com/hop-/goi/internal/core"
	"github.com/hop-/goi/internal/infra"
	"github.com/hop-/goi/internal/network"
	"github.com/hop-/golog"
)

func getCompressor(c *network.Connection) (compressors.Compressor, error) {
	// Read compressorType details
	_, b, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}

	compressorType := string(b)
	golog.Debug("Compression type", compressorType)

	compressor, err := compressors.New(compressorType)
	if err != nil {
		return compressor, err
	}

	err = c.WriteAll(network.OkRes)

	return compressor, err
}

func handleProducerHandshake(c *network.Connection) (*core.Producer, error) {
	// Send confirmation
	c.WriteAll(network.OkRes)

	// Read producer details
	_, producerName, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}

	producer, err := infra.NewProducer(string(producerName))
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
		defer infra.RemoveProducer(producer.Name)
	}
	if err != nil {
		return err
	}
	compressor, err := getCompressor(c)
	if err != nil {
		return err
	}

	golog.Infof("New producer %s accepted", producer.Name)

	// Producer main loop
producerMainLoop:
	for {
		golog.Debug("Waiting message from producer", producer.Name)
		messageType, b, err := c.ReadMessage()
		if err != nil {
			return err
		}
		golog.Debug("New message has been received from producer", producer.Name)

		switch messageType {
		// Handle special codes
		case network.SpecialCode:
			code := b[0]
			// Exit code
			if code == network.ExitCode {
				golog.Info("Received exit code from producer", producer.Name)
				break producerMainLoop
			}
		// Ping pong health check
		case network.PingMessage:
			golog.Debug("Received ping from producer", producer.Name)
			continue producerMainLoop
		// Handle other
		case network.GeneralMessage:
			buff, err := compressor.Decompress(b)
			if err != nil {
				golog.Error("Failed to decompress", err.Error())
				continue producerMainLoop
			}

			message, err := infra.NewMessageFromBuff(buff)
			if err != nil {
				golog.Error("Failed to create message from buffer", err.Error())
				// TODO: handle the error properly
				continue producerMainLoop
			}

			golog.Debugf("New message on %s topic with length of %d from producer %s", message.Topic, len(message.Content), producer.Name)

			infra.ProcessNewMessage(message)

			err = c.WriteSpecialCode(network.GoodResCode)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
