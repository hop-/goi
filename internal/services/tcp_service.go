package services

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/hop-/goi/internal/handlers"
	"github.com/hop-/golog"
)

type TcpService struct {
	port      int
	isRunning bool
	tlsConf   *tls.Config
	listener  net.Listener
}

func NewTcpService(port int, certFile string, keyFile string) (*TcpService, error) {
	var tlsConf *tls.Config

	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}
		tlsConf = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}
	return &TcpService{
		port:      port,
		isRunning: false,
		tlsConf:   tlsConf,
	}, nil
}

func (s *TcpService) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	var err error
	if s.tlsConf == nil {
		golog.Info("Starting TCP service")
		s.listener, err = net.Listen("tcp", addr)
	} else {
		golog.Info("Starting TCP+TLS service")
		s.listener, err = tls.Listen("tcp", addr, s.tlsConf)
	}
	if err != nil {
		return err
	}
	golog.Info("Listening TCP on", s.port)

	s.isRunning = true
	for s.isRunning {
		c, err := s.listener.Accept()
		if err != nil {
			golog.Error("Failed to accept TCP connection", err.Error())
			continue
		}
		golog.Info("New TCP connection accepted", c.RemoteAddr())

		// Each connection is handeled in separate goroutine
		go func() {
			// net.Conn implements network.SimpleConnection interface
			// no need for convertion

			err = handlers.ConnectionHandler(c)
			if err != nil {
				golog.Error("Failed to handle the connection", err.Error())
			}
		}()
	}

	return nil
}

func (s *TcpService) Stop() {
	s.isRunning = false

	s.listener.Close()
}
