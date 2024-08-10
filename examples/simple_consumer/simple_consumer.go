package main

import (
	"fmt"

	goi "github.com/hop-/goi/pkg"
)

func main() {
	topic := "test"
	config := goi.ConsumerConfig{Topic: &topic}

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
	fmt.Println("Reading a message")

	topic, message, err := c.ReadMessage()
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Message received on %s topic with message length of %d\n", topic, len(message))
}
