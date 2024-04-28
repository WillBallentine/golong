package broker

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type Broker struct {
	queues map[string]*Queue
	mu     sync.Mutex
}

type Queue struct {
	name     string
	messages []Message
	mu       sync.Mutex
	clients  map[net.Conn]struct{}
}

type Message struct {
	ID      string
	Payload []byte
}

func NewBroker() *Broker {
	broker := &Broker{
		queues: make(map[string]*Queue),
	}
	broker.NewQueue("init")
	broker.NewQueue("test")
	fmt.Println("new broker created")
	return broker
}

func (b *Broker) NewQueue(qName string) {
	// Send confirmation message back to client
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.queues[qName]; !ok {
		b.queues[qName] = &Queue{
			name:     qName,
			messages: make([]Message, 0),
		}
		fmt.Println("new queue created")
	}
}

func (b *Broker) Publish(message Message, name string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if queue, ok := b.queues[name]; ok {
		queue.messages = append(queue.messages, message)
		fmt.Println("Message added to queue:", name)
		//here is where we send new messages to current subscribers.
		//TODO: need to add a way to process messages added to queue when subscriber is not connected i.e. messages added before a subscriber connects
		queue.mu.Lock()
		for client := range queue.clients {
			fmt.Fprintf(client, "Message from server: %s\n", message.Payload)
		}
		queue.mu.Unlock()
	} else {
		fmt.Println("Queue not found:", name)
	}
}

// this is not used for now. Eventually, this will need to be implemented on the consumer side of things.
func (b *Broker) Consume(name string) (Message, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if queue, ok := b.queues[name]; ok {
		if len(queue.messages) == 0 {
			return Message{}, errors.New("queue is empty")
		}

		message := queue.messages[0]
		queue.messages = queue.messages[1:]

		return message, nil
	}

	return Message{}, errors.New("queue not found")
}

func (b *Broker) QueueHistory(name string) ([]Message, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if queue, ok := b.queues[name]; ok {
		return queue.messages, nil
	}

	return nil, errors.New("queue not found")
}

func (b *Broker) UserQueueCreate(qName string) {
	b.NewQueue(qName)
}

func (b *Broker) ProduceMessage(payload []byte, name string) {
	messageID := fmt.Sprintf("message-%d", time.Now().UnixNano())
	message := Message{ID: messageID, Payload: payload}

	b.Publish(message, name)
}

func (b *Broker) HandleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// Read incoming message
		message := scanner.Text()
		fmt.Println("Received message:", message)

		if strings.HasPrefix(message, "sub:") {
			queueName := strings.TrimPrefix(message, "sub:")
			b.subscribe(queueName, conn)
		} else {
			// Example: Split message into queue name and payload
			parts := strings.Split(message, ":")
			//there is a better way to do this....
			if len(parts) < 2 || len(parts) >= 4 {
				fmt.Println("Invalid message format")
				continue
			}

			if len(parts) == 2 {
				b.ProduceMessage([]byte(parts[1]), parts[0])
				// Send confirmation message back to client
				confirmation := "Message received and processed successfully"
				if _, err := conn.Write([]byte(confirmation + "\n")); err != nil {
					fmt.Println("Failed to send confirmation:", err)
					return
				}

			}

			if len(parts) == 3 {
				switch parts[2] {
				case "hist":
					history, err := b.QueueHistory(parts[0])
					if err != nil {
						fmt.Println("error retreiving history: ", err)
					}
					for _, value := range history {
						if _, err2 := conn.Write([]byte(string(value.Payload) + "\n")); err != nil {
							fmt.Println("Failed to send history: ", err2)
						}
					}
				case "nq":
					b.UserQueueCreate(parts[0])
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading from connection:", err)
		}
	}
}

func (b *Broker) subscribe(queueName string, conn net.Conn) {
	b.mu.Lock()
	defer b.mu.Unlock()

	q, ok := b.queues[queueName]
	if !ok {
		q = &Queue{
			clients: make(map[net.Conn]struct{}),
		}
		b.queues[queueName] = q
	}

	q.mu.Lock()
	defer q.mu.Unlock()
	if q.clients == nil {
		q.clients = make(map[net.Conn]struct{})
	}

	q.clients[conn] = struct{}{}
}
