package handlers

import (
	"fmt"

	"github.com/hop-/goi/internal/network"
)

func ConnectionHandler(c network.SimpleConnection) error {
	conn := network.NewConnection(c)

	// Get connection type
	connTypeBuff := make([]byte, 1)
	conn.ReadAll(connTypeBuff)
	connType := connTypeBuff[0]

	switch connType {
	case network.ConsumerType:
		return consumerHandler(conn)
	case network.ProducerType:
		return producerHandler(conn)
	default:
		return fmt.Errorf("conn: unknown connection type %b", connType)
	}
}

func consumerHandler(c *network.Connection) error {
	c.WriteAll(network.OkRes)

	// TODO: handle consumer
	return nil
}

func producerHandler(c *network.Connection) error {
	c.WriteAll(network.OkRes)

	// TODO: handle producer
	return nil
}
