package broker

import (
	"errors"
	"sync"
)

type Broker struct {
	queue Queue
	mu    sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{
		queue: Queue{},
	}
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
