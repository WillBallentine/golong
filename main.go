package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Producer struct {
	Broker *Broker
}

type Consumer struct {
	Broker *Broker
}

type Broker struct {
	queue Queue
	mu    sync.Mutex
}

type Queue struct {
	messages []Message
}

type Message struct {
	ID      string // message id
	Payload []byte // data of the message
	// possibly adding more fields later once implementation is fleshed out
}

func NewProducer(broker *Broker) *Producer {
	return &Producer{Broker: broker}
}

func NewBroker() *Broker {
	return &Broker{
		queue: Queue{},
	}
}

func (p *Producer) ProduceMessage(payload []byte) {
	messageID := fmt.Sprintf("message - %d", time.Now().UnixNano())
	message := Message{ID: messageID, Payload: payload}

	//TODO how do I publish to the broker?
	p.Broker.Publish(message)
	fmt.Sprintf("Message ready for publishing: %s $s", message.ID, message.Payload)

	fmt.Printf("Produced message with ID: %s\n", message.ID)
}

func (b *Broker) Publish(message Message) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.queue.messages = append(b.queue.messages, message)
}

func (b *Broker) Consume() (Message, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.queue.messages) == 0 {
		return Message{}, errors.New("queue is empty")
	}

	message := b.queue.messages[0]
	b.queue.messages = b.queue.messages[1:]
	return message, nil
}

func NewConsumer(broker *Broker) *Consumer {
	return &Consumer{Broker: broker}
}

func main() {
	broker := NewBroker()

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
