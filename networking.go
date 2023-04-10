package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

// The Server allows the transfer of data to a number of clients.
type Server struct {
	// clientMessagesBuffer is the number of messages that can be buffered
	// per client.
	//
	// The default amount is 32.
	clientMessagesBuffer int

	// publishRate is the rate at which the server can publish messages.
	//
	// The default is one publish every 50ms with a burst of 16.
	publishRate *rate.Limiter

	// logLocation controls where the server logs to.
	//
	// The default is to log to stdout.
	logLocation func(f string, v ...interface{})

	// serverMux is the http.ServeMuxer used to serve the server.
	serverMux http.ServeMux

	// clientsMu is a mutex used to synchronize access to the subscribers map.
	clientsMu sync.Mutex

	// clients is a map of all the clients connected to the server.
	clients map[*Client]struct{}
}

func newServer() *Server {
	server := &Server{
		clientMessagesBuffer: 32,
		clients:              make(map[*Client]struct{}),
		publishRate:          rate.NewLimiter(rate.Every(50*time.Millisecond), 16),
		logLocation:          log.Printf,
	}

	server.serverMux.Handle("/", http.FileServer(http.Dir(".")))
	server.serverMux.HandleFunc("/ws", server.subscriberHandler)
	server.serverMux.HandleFunc("/publish", server.publishHandler)

	return server
}

// clients represents a client connected to the server.
// If the client can't keep up with the server, the server will drop the client using closeSlow.
type Client struct {
	// closeSlow is a function that closes the client if it can't keep up with the server.
	closeSlow func()
	// msgs is a channel that the server sends messages to.
	msgs chan []byte
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.serverMux.ServeHTTP(w, r)
}

// subscriberHandler handles the /ws endpoint.
func (server *Server) subscriberHandler(w http.ResponseWriter, r *http.Request) {
	client, err := websocket.Accept(w, r, nil)
	if err != nil {
		server.logLocation("failed to accept websocket connection: %v", err)
		return
	}
	defer client.Close(websocket.StatusInternalError, "")

	// Subscribe a new client to the server.
	err = server.subscribe(r.Context(), client)
	if errors.Is(err, context.Canceled) {
		server.logLocation("subscriber disconnected")
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		server.logLocation("subscriber disconnected")
		return
	}

	if err != nil {
		server.logLocation("failed to subscribe: %v", err)
		return
	}
}

// publishHandler handles the /publish endpoint.
func (server *Server) publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body := http.MaxBytesReader(w, r.Body, 8192)
	msg, err := ioutil.ReadAll(body)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusRequestEntityTooLarge),
			http.StatusRequestEntityTooLarge,
		)
		return
	}
	fmt.Println(string(msg))

	server.publish(msg)

	w.WriteHeader(http.StatusAccepted)
}

// subscribe subscribes a client to the server.
// It creates a subscriber with a buffered msgs chan to give some room to slower
// connections and then registers the subscriber. It then listens for all messages
// and writes them to the WebSocket. If the context is cancelled or
// an error occurs, it returns and deletes the subscription.
//
// It uses CloseRead to keep reading from the connection to process control
// messages and cancel the context if the connection drops.
func (server *Server) subscribe(ctx context.Context, c *websocket.Conn) error {
	ctx = c.CloseRead(ctx)

	client := &Client{
		msgs: make(chan []byte, server.clientMessagesBuffer),

		closeSlow: func() {
			c.Close(websocket.StatusPolicyViolation, "too slow")
		},
	}

	server.addClient(client)
	defer server.removeClient(client)

	for {
		select {
		case msg := <-client.msgs:
			err := writeTimout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// publish publishes a message to all the clients.
func (server *Server) publish(msg []byte) {
	server.clientsMu.Lock()
	defer server.clientsMu.Unlock()

	for client := range server.clients {
		select {
		case client.msgs <- msg:
		default:
			client.closeSlow()
		}
	}
}

// addClient adds a client to the server.
func (server *Server) addClient(client *Client) {
	server.clientsMu.Lock()

	server.clients[client] = struct{}{}
	server.clientsMu.Unlock()
}

// removeClient removes a client from the server.
func (server *Server) removeClient(client *Client) {
	server.clientsMu.Lock()

	delete(server.clients, client)
	server.clientsMu.Unlock()
}

// writeTimout writes a message to the client with a timeout.
func writeTimout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}
