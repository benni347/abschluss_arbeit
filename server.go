package main

import (
	"fmt"
	"io"
	"net"
)

func printError(message string, err error) {
	fmt.Printf("ERROR: %s: %v\n", message, err)
}

func printInfo(message string, verbose bool) {
	if verbose {
		fmt.Printf("\033[1m%s\033[0m: %s\n", "INFO", message)
	}
}

func main() {
	printInfo("Starting server", true)

	port := ":14121"
	printInfo("Listening on port "+port, true)

	ln, err := net.Listen("tcp", port)
	if err != nil {
		printError("Failed to listen", err)
		return
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			printError("Failed to accept connection", err)
			continue
		}

		printInfo("Accepted connection", true)

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	printInfo("Handling connection", true)

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				printInfo("Client disconnected", true)
			} else {
				printError("Failed to read from connection", err)
			}
			return
		}

		printInfo("Received: "+string(buf[:n]), true)

		_, err = conn.Write(buf[:n])
		if err != nil {
			printError("Failed to write to connection", err)
			return
		}
	}
}
