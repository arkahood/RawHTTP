package main

import (
	"RAWHTTP/internal/request"
	"RAWHTTP/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 8080

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	// Check if the request target is /yourproblem
	if req.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{
			StatusCode: 400,
			Message:    "Your problem is not my problem\n",
		}
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{
			StatusCode: 500,
			Message:    "My problem is not my problem\n",
		}
	}

	if req.RequestLine.RequestTarget == "/route1" {
		message := "This is Route1\n"
		_, err := w.Write([]byte(message))
		if err != nil {
			return &server.HandlerError{
				StatusCode: 500,
				Message:    "Internal Server Error: failed to write response",
			}
		}
		return nil
	}

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

func main() {
	server, err := server.Serve(handler, port)
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
