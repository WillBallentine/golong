package broker

import (
	"errors"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
)

// TODO: need to break out server and client into different packages
// TODO: write a more uniform TUI for controlls
// TODO: debug why multiple versions of a queue are created at times via the /nq or /sq command (see phone screenshot for details)
var (
	sessions       map[*ssh.Session]*Broker
	newQueueCmd    = regexp.MustCompile(`^/nq.*`)
	switchQueueCmd = regexp.MustCompile(`^/sq.*`)
	helpCmd        = regexp.MustCompile(`^/help.*`)
	exitCmd        = regexp.MustCompile(`^/exit.*`)
	histCmd        = regexp.MustCompile(`^/h.*`)
)

type Broker struct {
	brokerId int
	queues   []Queue
	mu       sync.Mutex
}

type Producer struct {
	Broker *Broker
}

// TODO: remove consumer logic here. building to not require a broker
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

// TODO: remove consumer logic here
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

// TODO: remove consumer logic here.
func (b *Broker) Consume(name string) (Message, error) {
	for i := 0; i < len(b.queues); i++ {
		if name == b.queues[i].name {
			b.mu.Lock()
			defer b.mu.Unlock()
			if len(b.queues[i].messages) == 0 {
				return Message{}, errors.New("queue is empty\n")
			}

			message := b.queues[i].messages[0]

			// we only save the last 10 messages in any one queue
			if len(b.queues[i].messages) > 10 {
				b.queues[i].messages = b.queues[i].messages[1:]
			}
			return message, nil
		}
	}
	return Message{Payload: []byte("queue not found")}, nil
}

func (b *Broker) QueueHistory(currentQueue string, newTerm *term.Terminal) {
	if len(b.queues) > 0 {
		for i := 0; i < len(b.queues); i++ {
			if b.queues[i].name == currentQueue {
				fmt.Printf("queue name: %s\n", b.queues[i].name)
				if len(b.queues[i].messages) > 0 {
					for m := 0; m < len(b.queues[i].messages); m++ {
						newTerm.Write([]byte(b.queues[i].messages[m].Payload))
						newTerm.Write([]byte("\n"))
					}
				}
			}
		}
	}

}

func (b *Broker) UserQueueCreate(term *term.Terminal, line string, producer *Producer) string {
	term.Write([]byte("enter new queue name: "))
	currentQueue := ""
	name, err := term.ReadLine()
	if err != nil {
		fmt.Printf("error creating queue: %s", err)
	}
	currentQueue = name
	fmt.Println(currentQueue + "\n")
	producer.ProduceMessage([]byte(name), currentQueue)
	return currentQueue
}

func (b *Broker) SwitchQueue(term *term.Terminal) string {
	currentQueue := ""
	term.Write([]byte("select a queue: "))
	if len(b.queues) > 0 {
		for _, q := range b.queues {
			term.Write([]byte(q.name))
			term.Write([]byte("\n"))
		}
		name, err := term.ReadLine()
		if err != nil {
			fmt.Printf("error switching queues: %s", err)
		}
		for a := 0; a < len(b.queues); a++ {
			if b.queues[a].name == name {
				currentQueue = name
				term.Write([]byte("switched to queue -- "))
				term.Write([]byte(currentQueue))
				fmt.Printf("switched to queue %s", currentQueue)
				return currentQueue
			}
		}
		term.Write([]byte("no valid queue"))
		fmt.Println("no valid queues")
	}
	return currentQueue
}

func (b *Broker) SessionManager(sess ssh.Session) {
	newTerm := term.NewTerminal(sess, fmt.Sprint("> "))
	producer := NewProducer(b)
	newTerm.Write([]byte("control terminal starting..."))
	newTerm.Write([]byte("\n"))
	newTerm.Write([]byte(`
 _______  _______  _        _______  _        _______ 
(  ____ \(  ___  )( \      (  ___  )( (    /|(  ____ \
| (    \/| (   ) || (      | (   ) ||  \  ( || (    \/
| |      | |   | || |      | |   | ||   \ | || |      
| | ____ | |   | || |      | |   | || (\ \) || | ____ 
| | \_  )| |   | || |      | |   | || | \   || | \_  )
| (___) || (___) || (____/\| (___) || )  \  || (___) |
(_______)(_______)(_______/(_______)|/    )_)(_______)
	`))
	newTerm.Write([]byte("\n"))
	currentQueue := "init"
	b.NewQueue(currentQueue)
	for {
		newTerm.Write([]byte("current queue: "))
		newTerm.Write([]byte(currentQueue))
		newTerm.Write([]byte(" --"))
		fmt.Printf("currentQueue = %s\n", currentQueue)
		line, err := newTerm.ReadLine()
		if err != nil {
			break
		}

		if len(line) > 0 {
			if string(line[0]) == "/" {
				switch {
				//TODO: need to add in queue delete function
				//TODO: need to add in session management
				//TODO: need to add server side logging
				case exitCmd.MatchString(string(line)):
					return
				case histCmd.MatchString(string(line)):
					b.QueueHistory(currentQueue, newTerm)
				case newQueueCmd.MatchString(string(line)):
					currentQueue = b.UserQueueCreate(newTerm, line, producer)
				case switchQueueCmd.MatchString(string(line)):
					prevQueue := currentQueue
					currentQueue = b.SwitchQueue(newTerm)
					if currentQueue == "" {
						currentQueue = prevQueue
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
