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
	GoodResCode byte = 'K'
	BadResCode  byte = 'N'
	ExitCode    byte = 'X'
	OkRes            = []byte{GoodResCode}

	// Special messages
	MessageRequest = "+"
)
