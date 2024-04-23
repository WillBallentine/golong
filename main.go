package main

import (
	"fmt"
	"log"

	"github.com/WillBallentine/golong/broker"
	"github.com/gliderlabs/ssh"
)

func main() {
	b := broker.NewBroker()
	ssh.Handle(func(s ssh.Session) {
		b.SessionManager(s)
	})

	log.Println("starting ssh server on prot 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))

	fmt.Println("past fatal")
	//consumer := broker.NewConsumer(b)
	consumedMessage, err := b.Consume()
	if err != nil {
		fmt.Println("error encountered: %d", err)
		panic("oops! message could not be consumed. please try again")
	}

	log.Println(string(consumedMessage.Payload))

}

func Start() {
	//broker := broker.NewBroker()

	//producer := NewProducer(broker)
	//producer.ProduceMessage([]byte("hello, world!"))

	//consumer := NewConsumer(broker)
	//consumedMessage, err := broker.Consume()
	//if err != nil {
	//fmt.Println("error encountered: %d", err)
	//panic("oops! message could not be consumed. please try again.")
	//}

	//println(string(consumedMessage.Payload))
}
