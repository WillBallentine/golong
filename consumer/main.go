package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:2222")
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		os.Exit(1)
	}
	fmt.Println("Connection established")
	defer conn.Close()

	if _, err := conn.Write([]byte("sub:test\n")); err != nil {
		fmt.Println("Failed to send queue name:", err)
		os.Exit(1)
	}

	// Start listening for messages from the server
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// Read the message from the server
		message := scanner.Text()

		// Print the received message
		fmt.Println("Message from server:", message)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from connection:", err)
		os.Exit(1)
	}
}
