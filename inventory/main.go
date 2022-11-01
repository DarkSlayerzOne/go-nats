package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {

	wait := make(chan bool)

	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatal(err)
	}
	const topic = "my_order"

	nc.Subscribe(topic, func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
		nc.Publish(m.Reply, []byte("Hello"))
	})

	fmt.Println("Subscribed to", topic)

	<-wait

}
