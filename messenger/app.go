package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/benni347/encryption"
	"nhooyr.io/websocket"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	done := make(chan bool)
	go run(done)
}

func run(done chan bool) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("listening on http://%v", l.Addr())

	server := newServer()
	s := &http.Server{
		Handler:      server,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	log.Fatal(s.Shutdown(ctx))
	done <- false
}

func (a *App) Publisher(msg string, urlStr string) error {
	encryption.PrintInfo(fmt.Sprintf("Publishing message: %s", msg), true)
	encryption.PrintInfo(fmt.Sprintf("Publishing to: %s", urlStr), true)

	// Parse the WebSocket URL.
	u, err := url.Parse("ws://" + urlStr + "/publish")
	if err != nil {
		encryption.PrintError("Error during parsing url: ", err)
		return fmt.Errorf("Publish Failed: %s", err)
	}

	// Dial the WebSocket server.
	ctx := context.Background()
	encryption.PrintInfo("Dialing WebSocket...", true)
	encryption.PrintInfo(fmt.Sprintf("The websocket ctx is: %s", ctx), true)
	conn, _, err := websocket.Dial(ctx, u.String(), nil)
	if err != nil {
		encryption.PrintError("Error during dialing: ", err)
		return fmt.Errorf("Publish Failed: %s", err)
	}
	defer conn.Close(websocket.StatusInternalError, "WebSocket connection closed")

	// Send the message over the WebSocket connection.
	err = conn.Write(ctx, websocket.MessageText, []byte(msg))
	if err != nil {
		encryption.PrintError("Error during writing: ", err)
		return fmt.Errorf("Publish Failed: %s", err)
	}

	return nil
}

func (a *App) Dial(location string) (string, error) {
	encryption.PrintInfo(fmt.Sprintf("Dialing: %s", location), true)
	u := url.URL{Scheme: "ws", Host: location, Path: "/ws"}

	conn, _, err := websocket.Dial(context.Background(), u.String(), nil)
	if err != nil {
		encryption.PrintError("Failed to dial WebSocket: ", err)
		return "", fmt.Errorf("Dial Failed: %s", err)
	}
	defer conn.Close(websocket.StatusInternalError, "WebSocket connection closed")

	// Handle WebSocket events.
	encryption.PrintInfo("WebSocket connected", true)
	for {
		messageType, message, err := conn.Read(context.Background())
		if err != nil {
			encryption.PrintError("Failed to read WebSocket message: ", err)
			break
		}

		if messageType != websocket.MessageText {
			fmt.Println("Unexpected message type:", messageType)
			continue
		}

		fmt.Println("Received message from WebSocket:", string(message))
		return string(message), nil
	}
	return "", nil
}
