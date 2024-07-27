package core

import (
	"fmt"
	"sync"
)

type Storage interface {
	Init() error
	Topics() ([]Topic, error)
	NewTopic(Topic) error
	ConsumerGroups() ([]ConsumerGroup, error)
	NewConsumerGroup(ConsumerGroup) error
	Messages(Topic) ([]Message, error)
	NewMessage(Message) error
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
func (asc *AtomicStorageContainer) Topics() ([]Topic, error) {
	asc.mu.Lock()
	defer asc.mu.Unlock()

	return asc.storage.Topics()
}
func (asc *AtomicStorageContainer) NewTopic(t Topic) error {
	asc.mu.Lock()
	defer asc.mu.Unlock()

	return asc.storage.NewTopic(t)
}
func (asc *AtomicStorageContainer) ConsumerGroups() ([]ConsumerGroup, error) {
	asc.mu.Lock()
	defer asc.mu.Unlock()

	return asc.storage.ConsumerGroups()
}
func (asc *AtomicStorageContainer) NewConsumerGroup(cg ConsumerGroup) error {
	asc.mu.Lock()
	defer asc.mu.Unlock()

	return asc.storage.NewConsumerGroup(cg)
}
func (asc *AtomicStorageContainer) Messages(t Topic) ([]Message, error) {
	asc.mu.Lock()
	defer asc.mu.Unlock()

	return asc.storage.Messages(t)
}
func (asc *AtomicStorageContainer) NewMessage(m Message) error {
	asc.mu.Lock()
	defer asc.mu.Unlock()

	return asc.storage.NewMessage(m)
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
