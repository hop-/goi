package storages

import (
	"fmt"
	"sync"
)

type Storage interface {
	Close()
	// TODO: add methods
}

type StorageGenerator func(uri string) (Storage, error)

var (
	storageGenerators = make(map[string]StorageGenerator)
	storageInstance   Storage
	storageMutex      = &sync.Mutex{}
)

// Run init when package is imported
func init() {
	storageMutex.Lock()
	defer storageMutex.Unlock()

	storageGenerators["sqlite"] = newSqliteStorage
	// TODO: add all storage generators here
}

func InitStorage(storageType string, uri string) error {
	if storageInstance == nil {
		storageMutex.Lock()
		defer storageMutex.Unlock()

		if storageInstance == nil {
			generator, ok := storageGenerators[storageType]
			if !ok {
				return fmt.Errorf("unknown storage type: %s", storageType)
			}

			var err error
			storageInstance, err = generator(uri)

			return err
		}
	}

	return fmt.Errorf("storage should be initialize once")
}

func GetStorage() Storage {
	return storageInstance
}
