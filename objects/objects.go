package objects

import "sync"

type Message struct {
	ID      string // message id
	Payload []byte // data of the message
	// possibly adding more fields later once implementation is fleshed out
}

type Queue struct {
	messages []Message
}

type Broker struct {
	queue Queue
	mu    sync.Mutex
}
