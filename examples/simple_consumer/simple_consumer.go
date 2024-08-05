package main

import (
	"fmt"

	goi "github.com/hop-/goi/pkg"
)

func main() {
	config := goi.ConsumerConfig{}

	c, err := goi.NewConsumer(config)
	if err != nil {
		panic(err.Error())
	}

	err = c.Connect()
	if err != nil {
		panic(err.Error())
	}
	defer c.Disconnect()

	fmt.Println("Consumer has been connected")
}
