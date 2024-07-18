package network

import "encoding/binary"

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

func (c *Connection) ReadMessage() ([]byte, error) {
	// Read the message size
	var messageSize int64
	binary.Read(c, binary.LittleEndian, &messageSize)

	m := make([]byte, messageSize)
	offset := 0

	// Read whole message
	for int64(offset) < messageSize {
		size, err := c.conn.Read(m[offset:])
		if err != nil {
			return nil, err
		}

		offset += size
		// TODO: check size == 0 case
	}

	return m, nil
}

func (c *Connection) WriteMessage(m []byte) error {
	// Write the message size
	messageSize := len(m)
	binary.Write(c, binary.LittleEndian, int64(messageSize))

	offset := 0

	// Write whole message
	for offset < messageSize {
		size, err := c.conn.Write(m[offset:])
		if err != nil {
			return err
		}

		offset += size
		// TODO: check size == 0 case
	}

	return nil
}
