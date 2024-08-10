package core

type Producer struct {
	Name string
}

func NewProducer(name string) *Producer {
	return &Producer{
		Name: name,
	}
}
