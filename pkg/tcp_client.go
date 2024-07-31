package goi

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/hop-/goi/internal/network"
)

func connectTls(host string, port int) (network.SimpleConnection, error) {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	return tls.Dial("tcp", addr, config)
}

func connectTcp(host string, port int) (network.SimpleConnection, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	return net.Dial("tcp", addr)
}
