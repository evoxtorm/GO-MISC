package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Message struct {
	Body map[string]interface{} `json:"body"`
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Error while connecting to the server : %s\n", err.Error())
	}

	message := Message{
		Body: make(map[string]interface{}),
	}
	message.Body["v"] = "some value which should be seen"
	// message := "Hello to the server, can you able to recieve this message"

	jsondata, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Error while marshaling JSON: %s\n", err.Error())
		return
	}
	val, err := conn.Write([]byte(jsondata))
	fmt.Printf("%+v\n", val)
	if err != nil {
		fmt.Printf("Error while writing message to the server : %s\n", err.Error())
		return
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading from server: %s\n", err.Error())
		return
	}
	var newMessage Message
	err = json.Unmarshal(buffer[:n], &newMessage)
	if err != nil {
		fmt.Printf("Error while reading the json %s\n", err.Error())
	}
	body := newMessage.Body

	fmt.Println("Received Body:")
	for key, value := range body {
		fmt.Printf("%s: %v\n", key, value)
	}
	conn.Close()
}
