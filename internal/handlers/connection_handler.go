package handlers

import "github.com/hop-/goi/internal/network"

func ConnectionHandler(c network.SimpleConnection) {
	_ = network.NewConnection(c)
	// TODO
}
