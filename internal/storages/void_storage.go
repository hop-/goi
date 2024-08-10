package storages

import "github.com/hop-/goi/internal/core"

type VoidStorage struct{}

func (s *VoidStorage) Init() error {
	return nil
}

func (s *VoidStorage) Close() error {
	return nil
}

func (s *VoidStorage) Topics() ([]core.Topic, error) {
	return []core.Topic{}, nil
}

func (s *VoidStorage) NewTopic(*core.Topic) error {
	return nil
}

func (s *VoidStorage) ConsumerGroups() ([]core.ConsumerGroup, error) {
	return []core.ConsumerGroup{}, nil
}

func (s *VoidStorage) NewConsumerGroup(*core.ConsumerGroup) error {
	return nil
}

func (s *VoidStorage) Messages(*core.Topic) ([]core.Message, error) {
	return []core.Message{}, nil
}

func (s *VoidStorage) NewMessage(*core.Message) error {
	return nil
}

// TODO: add Storage implementation

func newVoidStorage(string) (Storage, error) {
	return &VoidStorage{}, nil
}
