package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Message struct {
	Body map[string]interface{} `json:"body"`
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Connection closed by the client")
				return
			}
			fmt.Printf("Error reading from connection: %s\n", err.Error())
			return
		}
		var message Message
		err = json.Unmarshal(buffer[:n], &message)

		body := message.Body

		fmt.Println("Received Body:")
		for key, value := range body {
			fmt.Printf("%s: %v\n", key, value)
		}
		jsonData, err := json.Marshal(message)
		if err != nil {
			fmt.Printf("Error marshaling message to JSON: %s\n", err.Error())
			return
		}
		_, err = conn.Write([]byte(jsonData))
		if err != nil {
			fmt.Printf("Error writing to connection: %s\n", err.Error())
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Error while creating the Listener: %s\n", err.Error())
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error while accepting the connection: %s\n", err.Error())
			return
		}

		fmt.Println("Client connected:", conn.RemoteAddr())

		go handleConnection(conn)
	}
}
