package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func sendFileToTCPSocket(filePath string, address string) error {

	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(file)
	if err != nil {
		return err
	}

	fmt.Println("message sent successfully")
	return nil
}

func main() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	filePath := dirname + "/output.bin"
	serverAddress := "localhost:5015"

	err = sendFileToTCPSocket(filePath, serverAddress)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
