package core

import (
	"fmt"
	"sync"
)

type Producer struct {
	Name string
}

var (
	producers = make(map[string]*Producer)
	pMu       = sync.Mutex{}
)

func newProducer(name string) (*Producer, error) {
	p := &Producer{
		Name: name,
	}

	err := addProducer(p)

	return p, err
}

func addProducer(p *Producer) error {
	pMu.Lock()
	defer pMu.Unlock()

	if _, ok := producers[p.Name]; ok {
		return fmt.Errorf("producer with %s name is already regisered", p.Name)
	}

	producers[p.Name] = p

	return nil
}

func removeProducer(p *Producer) error {
	pMu.Lock()
	defer pMu.Unlock()

	delete(producers, p.Name)
	return nil
}

func findProducerByName(name string) *Producer {
	pMu.Lock()
	defer pMu.Unlock()

	if p, ok := producers[name]; ok {
		return p
	}

	return nil
}
