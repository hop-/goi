package goi

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/hop-/goi/internal/network"
	"github.com/quic-go/quic-go"
)

func connectQuic(host string, port int) (network.SimpleConnection, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := quic.DialAddr(context.Background(), addr, tlsConfig, &quic.Config{})
	if err != nil {
		return nil, err
	}

	return conn.OpenStream()
}
