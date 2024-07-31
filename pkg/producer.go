package goi

type ProducerConfig struct {
	Host        *string
	Port        *int
	TcpPort     *int
	TcpFallback *bool
}

type Producer struct {
}

func NewProducer(config ProducerConfig) *Producer {
	return &Producer{}
}

func (p *Producer) Connect() error {
	// TODO
	return nil
}

func (p *Producer) Disconnect() error {
	// TODO
	return nil
}

// TODO
