package broker

import (
	"errors"
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
	"sync"
	"time"
)

var (
	sessions map[*ssh.Session]*Broker
)

type Broker struct {
	brokerId int
	queue    Queue
	Users    []User
	mu       sync.Mutex
}

type User struct {
	Session  ssh.Session
	Terminal *term.Terminal
}
type Producer struct {
	Broker *Broker
}

type Consumer struct {
	Broker *Broker
}

func NewBroker() *Broker {
	return &Broker{
		queue: Queue{},
	}
}

func NewConsumer(broker *Broker) *Consumer {
	return &Consumer{Broker: broker}
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

func (b *Broker) SessionManager(sess ssh.Session) {
	newTerm := term.NewTerminal(sess, fmt.Sprint("> "))
	for {
		line, err := newTerm.ReadLine()
		if err != nil {
			break
		}

		if len(line) > 0 {
			producer := NewProducer(b)
			producer.ProduceMessage([]byte(line))
		}
	}
}
func NewProducer(broker *Broker) *Producer {
	return &Producer{Broker: broker}
}

func (p *Producer) ProduceMessage(payload []byte) {
	messageID := fmt.Sprintf("message - %d", time.Now().UnixNano())
	message := Message{ID: messageID, Payload: payload}

	p.Broker.Publish(message)
	fmt.Sprintf("Message ready for publishing: %s $s", message.ID, message.Payload)

	fmt.Printf("Produced message with ID: %s\n", message.ID)

	consumedMessage, err := p.Broker.Consume()
	if err != nil {
		fmt.Println("error encountered: %d", err)
		panic("oops! message could not be consumed. please try again")
	}

	fmt.Println(string(consumedMessage.Payload))

}
