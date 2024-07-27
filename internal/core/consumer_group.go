package core

type ConsumerGroup struct {
	Id   *int
	Name string
}

var (
	consumerGroups = make(map[string]ConsumerGroup)
)

func AddConsumerGroup(cg ConsumerGroup) error {
	err := GetStorage().NewConsumerGroup(cg)
	if err != nil {
		return err
	}

	// It should be good to go when storage didn't return an error
	consumerGroups[cg.Name] = cg

	return nil
}

func LoadConsumerGroups() error {
	cgs, err := GetStorage().ConsumerGroups()
	if err != nil {
		return err
	}

	for _, cg := range cgs {
		consumerGroups[cg.Name] = cg
	}

	return nil
}
