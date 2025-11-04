package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"RAWHTTP/internal/request"
	"RAWHTTP/pkg/server"
)

const port = 8080

func handler1(w io.Writer, req *request.Request) *server.HandlerError {
	// Write a simple response
	message := "Hello World!\n"
	_, err := w.Write([]byte(message))
	if err != nil {
		return &server.HandlerError{
			StatusCode: 500,
			Message:    "Internal Server Error: failed to write response",
		}
	}
	return nil
}

func handler2(w io.Writer, req *request.Request) *server.HandlerError {
	// Write a simple response
	message := "Hello World2!\n"
	_, err := w.Write([]byte(message))
	if err != nil {
		return &server.HandlerError{
			StatusCode: 500,
			Message:    "Internal Server Error: failed to write response",
		}
	}
	return nil
}

func main() {
	r := server.NewRouter()
	server, err := server.Serve(r, port)

	r.AddRoute("/", handler1)
	r.AddRoute("/second-route", handler2)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
