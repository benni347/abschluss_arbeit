package main

import (
	"fmt"
	"net"
)

func printError(message string, err error) {
	fmt.Printf("ERROR: %s: %v\n", message, err)
}

func printInfo(message string, verbose bool) {
	if verbose {
		fmt.Printf("INFO: %s\n", message)
	}
}

func main() {
	printInfo("Starting client", true)

	port := ":14121"
	addr := "localhost" + port

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		printError("Failed to listen", err)
		return
	}

	defer conn.Close()

	// Send "Hello World" to the server
	data := "Hello World"
	_, err = conn.Write([]byte(data))
	if err != nil {
		printError("Failed to write to server", err)
		return
	}

	printInfo("Sent "+data+" to server", true)

	// Read the server's response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		printError("Failed to read from server", err)
		return
	}

	response := string(buf[:n])
	printInfo("Received "+response+" from server", true)
}
