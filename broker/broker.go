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
	queues   []Queue
	mu       sync.Mutex
}

type Producer struct {
	Broker *Broker
}

type Consumer struct {
	Broker *Broker
}

func NewBroker() *Broker {
	return &Broker{
		queues: make([]Queue, 0),
	}
}

func (b *Broker) NewQueue(qName string) Queue {
	return Queue{
		name:     qName,
		messages: make([]Message, 0),
	}
}

func NewConsumer(broker *Broker) *Consumer {
	return &Consumer{Broker: broker}
}

func (b *Broker) Publish(message Message, name string) {
	for i := range b.queues {
		if name == b.queues[i].name {
			b.mu.Lock()
			defer b.mu.Unlock()
			b.queues[i].messages = append(b.queues[i].messages, message)
		}
	}
}

func (b *Broker) Consume(name string) (Message, error) {
	for i := range b.queues {
		if name == b.queues[i].name {
			b.mu.Lock()
			defer b.mu.Unlock()
			if len(b.queues[i].messages) == 0 {
				return Message{}, errors.New("queue is empty")
			}

			message := b.queues[i].messages[0]
			b.queues[i].messages = b.queues[i].messages[1:]
			return message, nil
		}
	}
	fmt.Printf("queue name: %s", name)
	return Message{Payload: []byte("queue not found")}, nil
}

func (b *Broker) SessionManager(sess ssh.Session) {
	newTerm := term.NewTerminal(sess, fmt.Sprint("> "))
	for {
		line, err := newTerm.ReadLine()
		if err != nil {
			break
		}

		if len(line) > 0 {
			//need to add in queue select
			//need to add in queue create
			producer := NewProducer(b)
			producer.ProduceMessage([]byte(line), "test")
		}
	}
}

func NewProducer(broker *Broker) *Producer {
	return &Producer{Broker: broker}
}

func (p *Producer) ProduceMessage(payload []byte, name string) {
	queueName := name
	fmt.Println("beginning produce")
	if len(p.Broker.queues) > 0 {
		for i := range p.Broker.queues {
			if p.Broker.queues[i].name == name {
				queueName = name
				fmt.Println("using existing queue")
				continue
			} else {
				fmt.Println("creating new queue")
				tmp := p.Broker.NewQueue(name)
				p.Broker.queues = append(p.Broker.queues, tmp)
				queueName = tmp.name
			}
		}
	} else {
		fmt.Println("produce else block")
		tmp := p.Broker.NewQueue(name)
		p.Broker.queues = append(p.Broker.queues, tmp)
		queueName = tmp.name
	}
	fmt.Println("creating message")
	messageID := fmt.Sprintf("message - %d", time.Now().UnixNano())
	message := Message{ID: messageID, Payload: payload}

	p.Broker.Publish(message, queueName)
	fmt.Sprintf("Message ready for publishing: %s $s", message.ID, message.Payload)

	fmt.Printf("Produced message with ID: %s\n", message.ID)

	consumedMessage, err := p.Broker.Consume(queueName)
	if err != nil {
		fmt.Println("error encountered: %d", err)
		panic("oops! message could not be consumed. please try again")
	}

	fmt.Println(string(consumedMessage.Payload))

}
