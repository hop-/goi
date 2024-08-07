package network

import (
	"encoding/binary"
)

type Connection struct {
	conn SimpleConnection
}

func NewConnection(conn SimpleConnection) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}

func (c *Connection) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *Connection) ReadAll(b []byte) error {
	offset := 0

	// Read whole message
	for offset < len(b) {
		size, err := c.conn.Read(b[offset:])
		if err != nil {
			return err
		}

		offset += size
		// TODO: check size == 0 case
	}

	return nil
}

func (c *Connection) WriteAll(b []byte) error {
	offset := 0

	// Write whole message
	for offset < len(b) {
		size, err := c.conn.Write(b[offset:])
		if err != nil {
			return err
		}

		offset += size
		// TODO: check size == 0 case
	}

	return nil
}

func (c *Connection) ReadMessage() (int, []byte, error) {
	// Read the message size
	var messageSize int64
	err := binary.Read(c, binary.LittleEndian, &messageSize)
	if err != nil {
		return 0, nil, err
	}

	messageType := GeneralMessage
	var m []byte

	switch messageSize {
	case SpecialCode:
		messageType = SpecialCode
		m = make([]byte, 1)
		// Read whole message
		err = c.ReadAll(m)
	case SpecialMessage:
		messageType = SpecialMessage
		_, m, err = c.ReadMessage()
	case PingMessage:
		messageType = PingMessage
	default:
		m = make([]byte, messageSize)
		// Read whole message
		err = c.ReadAll(m)
	}

	return messageType, m, err
}

func (c *Connection) WriteMessage(m []byte) error {
	// Write the message size
	messageSize := len(m)
	err := binary.Write(c, binary.LittleEndian, int64(messageSize))
	if err != nil {
		return err
	}

	// Write whole message
	return c.WriteAll(m)
}

func (c *Connection) WriteSpecialCode(code byte) error {
	// Write the special number for message size
	err := binary.Write(c, binary.LittleEndian, int64(SpecialCode))
	if err != nil {
		return err
	}

	// Write whole message
	return c.WriteAll([]byte{code})
}

func (c *Connection) WriteSpecialMessage(m []byte) error {
	// Write the special number for message size
	err := binary.Write(c, binary.LittleEndian, int64(SpecialMessage))
	if err != nil {
		return err
	}

	// Write the message
	return c.WriteMessage(m)
}

func (c *Connection) Ping() error {
	return binary.Write(c, binary.LittleEndian, int64(PingMessage))
}

func (c *Connection) Close() error {
	return c.conn.Close()
}
