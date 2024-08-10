package core

type Topic struct {
	Name string
}

func NewTopic(name string) *Topic {
	return &Topic{
		Name: name,
	}
}
