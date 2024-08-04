package network

var (
	// Connection types
	ConsumerType        byte = 'C'
	ProducerType        byte = 'P'
	ConsumerTypeMessage      = []byte{ConsumerType}
	ProducerTypeMessage      = []byte{ProducerType}

	// Response codes
	OkResCode byte = 'K'
	ExitCode  byte = 'X'
	OkRes          = []byte{OkResCode}
)
