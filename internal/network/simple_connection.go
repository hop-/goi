package network

type SimpleConnection interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}
