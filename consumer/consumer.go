package consumer

type Consumer struct {
}

func NewConsumer() *Consumer {
	return nil
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

			// we only save the last 10 messages in any one queue
			if len(b.queues[i].messages) > 10 {
				b.queues[i].messages = b.queues[i].messages[1:]
			}
			return message, nil
		}
	}
	return Message{Payload: []byte("queue not found")}, nil
}

//consumedMessage, err := p.Broker.Consume(queueName)
//if err != nil {
//	fmt.Printf("error encountered: %d", err)
//	panic("oops! message could not be consumed. please try again")
//}
//
//	fmt.Println(string(consumedMessage.Payload))
