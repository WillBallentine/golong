package broker

import (
	"errors"
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
	"regexp"
	"sync"
	"time"
)

var (
	sessions       map[*ssh.Session]*Broker
	newQueueCmd    = regexp.MustCompile(`^/nq.*`)
	switchQueueCmd = regexp.MustCompile(`^/sq.*`)
	helpCmd        = regexp.MustCompile(`^/help.*`)
	exitCmd        = regexp.MustCompile(`^/exit.*`)
	listCmd        = regexp.MustCompile(`^/list.*`)
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
	for i := 0; i < len(b.queues); i++ {
		if name == b.queues[i].name {
			b.mu.Lock()
			defer b.mu.Unlock()
			b.queues[i].messages = append(b.queues[i].messages, message)
			fmt.Println("message added to queue\n")
			break
		}
	}
}

func (b *Broker) Consume(name string) (Message, error) {
	for i := 0; i < len(b.queues); i++ {
		if name == b.queues[i].name {
			b.mu.Lock()
			defer b.mu.Unlock()
			if len(b.queues[i].messages) == 0 {
				return Message{}, errors.New("queue is empty\n")
			}

			message := b.queues[i].messages[0]
			b.queues[i].messages = b.queues[i].messages[1:]
			return message, nil
		}
	}
	return Message{Payload: []byte("queue not found")}, nil
}

func (b *Broker) SessionManager(sess ssh.Session) {
	newTerm := term.NewTerminal(sess, fmt.Sprint("> "))
	producer := NewProducer(b)
	currentQueue := "init"
	for {
		fmt.Printf("currentQueue = %s\n", currentQueue)
		line, err := newTerm.ReadLine()
		if err != nil {
			break
		}

		if len(line) > 0 {
			if string(line[0]) == "/" {
				switch {
				//TODO: need to add in queue history print
				case exitCmd.MatchString(string(line)):
					return
				case newQueueCmd.MatchString(string(line)):
					newTerm.Write([]byte("enter new queue name: "))
					name, err := newTerm.ReadLine()
					if err != nil {
						break
					}
					currentQueue = name
					producer.ProduceMessage([]byte(line), currentQueue)
					fmt.Println(currentQueue + "\n")
				case switchQueueCmd.MatchString(string(line)):
					newTerm.Write([]byte("select a queue: "))
					if len(b.queues) > 0 {
						for i := range b.queues {
							newTerm.Write([]byte("\n"))
							newTerm.Write([]byte(b.queues[i].name))
							newTerm.Write([]byte("\n"))
						}
						name, err := newTerm.ReadLine()
						if err != nil {
							break
						}
						currentQueue = name
						newTerm.Write([]byte("switched to queue -- "))
						newTerm.Write([]byte(currentQueue))
						fmt.Printf("switched to queue %s", currentQueue)
						break
					} else {
						fmt.Println("not a valid queue. please enter a new cmd")
						break
					}

				default:
					producer.ProduceMessage([]byte(line), currentQueue)
					fmt.Printf("default path. Queue = %s\n", currentQueue)
				}
			} else {
				producer.ProduceMessage([]byte(line), currentQueue)
			}
		}
	}
}

func NewProducer(broker *Broker) *Producer {
	return &Producer{Broker: broker}
}

func (p *Producer) ProduceMessage(payload []byte, name string) {
	queueName := name
	if len(p.Broker.queues) > 0 {
		for i := 0; i < len(p.Broker.queues); i++ {
			if p.Broker.queues[i].name == name {
				queueName = name
				break
			} else {
				tmp := p.Broker.NewQueue(name)
				p.Broker.queues = append(p.Broker.queues, tmp)
				queueName = tmp.name
			}
		}
	} else {
		tmp := p.Broker.NewQueue(name)
		p.Broker.queues = append(p.Broker.queues, tmp)
		queueName = tmp.name
	}
	messageID := fmt.Sprintf("message - %d", time.Now().UnixNano())
	message := Message{ID: messageID, Payload: payload}

	fmt.Printf("Message ready for publishing: %s", message.ID)
	p.Broker.Publish(message, queueName)

	fmt.Printf("Produced message with ID: %s\n", message.ID)

	//TODO: make consumer logic buildable for outside users. This is temp static consumer
	consumedMessage, err := p.Broker.Consume(queueName)
	if err != nil {
		fmt.Printf("error encountered: %d", err)
		panic("oops! message could not be consumed. please try again")
	}

	fmt.Println(string(consumedMessage.Payload))

}
