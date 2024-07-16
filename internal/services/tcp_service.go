package services

import (
	"crypto/tls"
	"fmt"
	"net"

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
		s.listener, err = net.Listen("tcp", addr)
	} else {
		s.listener, err = tls.Listen("tcp", addr, s.tlsConf)
	}
	if err != nil {
		return err
	}

	s.isRunning = true
	for s.isRunning {
		c, err := s.listener.Accept()
		if err != nil {
			golog.Error("Failed to accept TCP connection", err.Error())
		}

		// TODO
		_ = c
	}

	return nil
}

func (s *TcpService) Stop() {
	s.isRunning = false

	s.listener.Close()
}
