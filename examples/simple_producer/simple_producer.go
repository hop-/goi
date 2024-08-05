package main

import (
	"fmt"

	goi "github.com/hop-/goi/pkg"
)

func main() {
	config := goi.ProducerConfig{}

	p, err := goi.NewProducer(config)
	if err != nil {
		panic(err.Error())
	}

	err = p.Connect()
	if err != nil {
		panic(err.Error())
	}
	defer p.Disconnect()

	fmt.Println("Producer has been connected")
}
