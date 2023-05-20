package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"net"
)

type Message struct {
	Body map[string]interface{} `json:"body"`
}

func encrypt(message Message) ([]byte, error) {
	text, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}
	key := []byte("0123456789abcdef0123456789abcdef")
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}
	iv := make([]byte, aes.BlockSize)

	cipherText := make([]byte, len(text))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, text)
	return cipherText, nil
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
	encryptedData, err := encrypt(message)
	val, err := conn.Write(encryptedData)
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
