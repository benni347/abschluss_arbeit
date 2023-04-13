package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
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

func (a *App) Publisher(msg string, url string) {
	url = "http://" + url + "/publish"

}

func dial(url string) {
	conn := "ws://" + url + "/ws"
}
