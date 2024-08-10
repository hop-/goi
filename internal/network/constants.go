package network

const (
	GeneralMessage = 0
	SpecialCode    = -1
	SpecialMessage = -2
	PingMessage    = -3
)

var (
	// Connection types
	ConsumerType        byte = 'C'
	ProducerType        byte = 'P'
	ConsumerTypeMessage      = []byte{ConsumerType}
	ProducerTypeMessage      = []byte{ProducerType}

	// Response codes
	OkResCode  byte = 'K'
	BadResCode byte = 'N'
	ExitCode   byte = 'X'
	OkRes           = []byte{OkResCode}

	// Special messages
	MessageRequest = "+"
)
