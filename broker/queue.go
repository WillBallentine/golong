package broker

type Queue struct {
	name     string
	messages []Message
}

type Message struct {
	ID      string // message id
	Payload []byte // data of the message
	// possibly adding more fields later once implementation is fleshed out
	// future func: I want users to be able to establish a connection (tcp? ssh?), designate a queue, and send a message
}
