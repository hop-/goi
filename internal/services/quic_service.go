package services

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/hop-/golog"
	quic "github.com/quic-go/quic-go"
)

type QuicService struct {
	port      int
	tlsConf   *tls.Config
	quicConf  *quic.Config
	isRunning bool
	listener  *quic.Listener
}

func NewQuicService(port int, certFile string, keyFile string) (*QuicService, error) {
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

	return &QuicService{
		port:     port,
		tlsConf:  tlsConf,
		quicConf: &quic.Config{
			// TODO
		},
	}, nil
}

func (s *QuicService) Start() error {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: s.port})
	if err != nil {
		return err
	}
	defer conn.Close()

	tr := quic.Transport{
		Conn: conn,
	}
	defer tr.Close()

	s.listener, err = tr.Listen(s.tlsConf, s.quicConf)
	if err != nil {
		return err
	}

	s.isRunning = true
	for s.isRunning {
		c, err := s.listener.Accept(context.Background())
		if err != nil {
			golog.Error("Failed to accept QUIC connection", err.Error())
		}

		// TODO
		_ = c
	}

	return nil
}

func (s *QuicService) Stop() {
	s.isRunning = false

	if s.listener != nil {
		s.listener.Close()
	}
}
