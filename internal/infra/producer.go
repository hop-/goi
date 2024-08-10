package infra

import (
	"fmt"
	"sync"

	"github.com/hop-/goi/internal/core"
)

var (
	producers = make(map[string]*core.Producer)
	pMu       = sync.Mutex{}
)

func newProducer(name string) (*core.Producer, error) {
	p := core.NewProducer(name)

	err := addProducer(p)

	return p, err
}

func addProducer(p *core.Producer) error {
	pMu.Lock()
	defer pMu.Unlock()

	if _, ok := producers[p.Name]; ok {
		return fmt.Errorf("producer with %s name is already regisered", p.Name)
	}

	producers[p.Name] = p

	return nil
}

func removeProducer(p *core.Producer) error {
	pMu.Lock()
	defer pMu.Unlock()

	delete(producers, p.Name)
	return nil
}

func findProducerByName(name string) *core.Producer {
	pMu.Lock()
	defer pMu.Unlock()

	if p, ok := producers[name]; ok {
		return p
	}

	return nil
}
