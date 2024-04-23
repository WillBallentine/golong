package broker

type Queue struct {
	messages []Message
}

type Message struct {
	ID      string // message id
	Payload []byte // data of the message
	// possibly adding more fields later once implementation is fleshed out
}
