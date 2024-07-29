package handlers

import "github.com/hop-/goi/internal/network"

func producerHandler(c *network.Connection) error {
	// Send confirmation
	c.WriteAll(network.OkRes)

	// TODO: handle producer
	return nil
}
