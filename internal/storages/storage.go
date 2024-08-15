package storages

import (
	"fmt"
	"sync"

	"github.com/hop-/goi/internal/core"
)

type Storage interface {
	Init() error
	Topics() ([]core.Topic, error)
	NewTopic(*core.Topic) error
	ConsumerGroups() ([]core.ConsumerGroup, error)
	NewConsumerGroup(*core.ConsumerGroup) error
	Messages(*core.Topic) ([]core.Message, error)
	NewMessage(*core.Message) error
	NextMessageForConsumerGroup(*core.ConsumerGroup, *core.Topic) (*core.Message, error)
	Close() error
}

type AtomicStorageContainer struct {
	storage Storage
	mu      *sync.Mutex
}

func (asc *AtomicStorageContainer) Init() error {
	asc.mu.Lock()
	defer asc.mu.Unlock()

	return asc.storage.Init()
}

func (asc *AtomicStorageContainer) Topics() ([]core.Topic, error) {
	asc.mu.Lock()
	defer asc.mu.Unlock()

	return asc.storage.Topics()
}

func (asc *AtomicStorageContainer) NewTopic(t *core.Topic) error {
	asc.mu.Lock()
	defer asc.mu.Unlock()

	return asc.storage.NewTopic(t)
}

func (asc *AtomicStorageContainer) ConsumerGroups() ([]core.ConsumerGroup, error) {
	asc.mu.Lock()
	defer asc.mu.Unlock()

	return asc.storage.ConsumerGroups()
}

func (asc *AtomicStorageContainer) NewConsumerGroup(cg *core.ConsumerGroup) error {
	asc.mu.Lock()
	defer asc.mu.Unlock()

	return asc.storage.NewConsumerGroup(cg)
}

func (asc *AtomicStorageContainer) Messages(t *core.Topic) ([]core.Message, error) {
	asc.mu.Lock()
	defer asc.mu.Unlock()

	return asc.storage.Messages(t)
}

func (asc *AtomicStorageContainer) NewMessage(m *core.Message) error {
	asc.mu.Lock()
	defer asc.mu.Unlock()

	return asc.storage.NewMessage(m)
}

func (asc *AtomicStorageContainer) NextMessageForConsumerGroup(cg *core.ConsumerGroup, t *core.Topic) (*core.Message, error) {
	asc.mu.Lock()
	defer asc.mu.Unlock()

	return asc.storage.NextMessageForConsumerGroup(cg, t)
}

func (asc *AtomicStorageContainer) Close() error {
	asc.mu.Lock()
	defer asc.mu.Unlock()

	return asc.storage.Close()
}

type StorageGenerator func(uri string) (Storage, error)

var (
	storageGenerators        = make(map[string]StorageGenerator)
	storageInstanceContainer Storage
	storageMutex             = &sync.Mutex{}
)

func RegisterStorage(name string, generator StorageGenerator) {
	storageMutex.Lock()
	defer storageMutex.Unlock()

	storageGenerators[name] = generator
}

func InitStorage(storageType string, uri string) error {
	if storageInstanceContainer == nil {
		storageMutex.Lock()
		defer storageMutex.Unlock()

		if storageInstanceContainer == nil {
			generator, ok := storageGenerators[storageType]
			if !ok {
				return fmt.Errorf("unknown storage type: %s", storageType)
			}

			var err error
			storageInstance, err := generator(uri)
			if err != nil {
				return err
			}

			// Using atomic storage container instead of plane instance
			storageInstanceContainer = &AtomicStorageContainer{
				storage: storageInstance,
				mu:      &sync.Mutex{},
			}

			return storageInstanceContainer.Init()
		}
	}

	return fmt.Errorf("storage should be initialize once")
}

func GetStorage() Storage {
	return storageInstanceContainer
}
