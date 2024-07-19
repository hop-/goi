package network

var (
	// Connection types
	ConsumerType byte = 'C'
	ProducerType byte = 'P'

	// Response codes
	OkResCode byte = 'K'
	OkRes          = []byte{OkResCode}
)
