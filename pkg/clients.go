package goi

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/hop-/goi/internal/network"
	"github.com/quic-go/quic-go"
)

var (
	defaultHost     = "localhost"
	defaultPort     = 4554
	defaultFallback = false
)

func connectQuic(host string, port int) (network.SimpleConnection, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := quic.DialAddr(context.TODO(), addr, tlsConfig, &quic.Config{})
	if err != nil {
		return nil, err
	}

	return conn.OpenStream()
}
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

func connect(host string, quicPort int, tcpPort int, fallback bool) (network.SimpleConnection, error) {
	conn, err := connectQuic(host, quicPort)
	if err == nil || !fallback {
		return conn, err
	}

	conn, err = connectTls(host, tcpPort)
	if err == nil {
		return conn, nil
	}

	return connectTcp(host, tcpPort)
}
