package main

import (
	"log"

	"github.com/WillBallentine/golong/broker"
	"github.com/gliderlabs/ssh"
)

func main() {
	b := broker.NewBroker()
	ssh.Handle(func(s ssh.Session) {
		b.SessionManager(s)
	})

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))

}
