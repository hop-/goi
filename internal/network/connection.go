package network

import "encoding/binary"

const (
	GeneralMessage = 0
	SpecialCode    = 1
	SpecialMessage = 2
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
	binary.Read(c, binary.LittleEndian, &messageSize)

	messageType := GeneralMessage
	var m []byte
	var err error

	switch messageSize {
	case -1:
		m = make([]byte, 1)
		// Read whole message
		err = c.ReadAll(m)
	case -2:
		_, m, err = c.ReadMessage()
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
	binary.Write(c, binary.LittleEndian, int64(messageSize))

	// Write whole message
	return c.WriteAll(m)
}

func (c *Connection) WriteSpecialCode(code byte) error {
	// Write the special number for message size
	binary.Write(c, binary.LittleEndian, int64(-1))

	// Write whole message
	return c.WriteAll([]byte{code})
}

func (c *Connection) WriteSpecialMessage(m []byte) error {
	// Write the special number for message size
	binary.Write(c, binary.LittleEndian, int64(-2))

	// Write the message
	return c.WriteMessage(m)
}
