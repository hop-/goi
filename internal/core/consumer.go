package core

import "fmt"

type Consumer struct {
	Name  string
	group *ConsumerGroup
}

var (
	consumers = make(map[string]map[string]Consumer)
)

func AddConsumer(c Consumer) error {
	if cg, ok := consumers[c.group.Name]; !ok {
		if _, ok := cg[c.Name]; ok {
			return fmt.Errorf("Consumer with %s name already registered", c.Name)
		}

		cg[c.Name] = c

		return nil
	}

	consumers[c.group.Name] = map[string]Consumer{
		c.Name: c,
	}

	return nil
}
