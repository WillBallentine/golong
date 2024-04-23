package main

import (
	"fmt"
	"github.com/WillBallentine/golong/broker"
	"time"
)

type Producer struct {
	Broker *broker.Broker
}

type Consumer struct {
	Broker *broker.Broker
}

func NewProducer(broker *broker.Broker) *Producer {
	return &Producer{Broker: broker}
}

func (p *Producer) ProduceMessage(payload []byte) {
	messageID := fmt.Sprintf("message - %d", time.Now().UnixNano())
	message := broker.Message{ID: messageID, Payload: payload}

	//TODO how do I publish to the broker?
	p.Broker.Publish(message)
	fmt.Sprintf("Message ready for publishing: %s $s", message.ID, message.Payload)

	fmt.Printf("Produced message with ID: %s\n", message.ID)
}

func NewConsumer(broker *broker.Broker) *Consumer {
	return &Consumer{Broker: broker}
}

func main() {
	broker := broker.NewBroker()

	producer := NewProducer(broker)
	producer.ProduceMessage([]byte("hello, world!"))

	//consumer := NewConsumer(broker)
	consumedMessage, err := broker.Consume()
	if err != nil {
		fmt.Println("error encountered: %d", err)
		panic("oops!")
	}

	println(string(consumedMessage.Payload))
}
