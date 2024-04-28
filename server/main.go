package main

import (
	"fmt"
	"github.com/WillBallentine/golong/broker"
	"net"
)

func main() {
	b := broker.NewBroker()
	ln, err := net.Listen("tcp", "127.0.0.1:2222")
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on port 2222")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		go b.HandleConnection(conn)
	}
}
