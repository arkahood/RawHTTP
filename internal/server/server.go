package server

import (
	"RAWHTTP/internal/request"
	"RAWHTTP/internal/response"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
)

type Server struct {
	listener       net.Listener
	serverIsClosed atomic.Bool
	handler        Handler
}

type HandlerError struct {
	StatusCode int
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

// WriteHandlerError writes a HandlerError to an io.Writer with proper HTTP response format.
func (hErr *HandlerError) WriteHandlerError(w io.Writer) error {
	// Determine the status code
	var statusCode response.StatusCode
	switch hErr.StatusCode {
	case 400:
		statusCode = response.StatusCodeClientErrors
	case 500:
		statusCode = response.StatusCodeInternalServerError
	default:
		statusCode = response.StatusCodeInternalServerError
	}

	// Write status line
	err := response.WriteStatusLine(w, statusCode)
	if err != nil {
		return err
	}

	// Prepare error message as body
	body := []byte(hErr.Message)

	// Write headers with content length
	headers := response.GetDefaultHeaders(len(body))
	err = response.WriteHeaders(w, headers)
	if err != nil {
		return err
	}

	// Write body
	_, err = w.Write(body)
	return err
}

func Serve(h Handler, port int) (*Server, error) {
	lstnr, err := net.Listen("tcp", fmt.Sprint(":", port))
	if err != nil {
		return nil, errors.New("tcp error happened")
	}
	server := Server{
		listener: lstnr,
		handler:  h,
	}

	go server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	s.serverIsClosed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if !s.serverIsClosed.Load() {
				fmt.Println("Accept error:", err.Error())
			}
			return
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	// Parse the request from the connection
	req, err := request.RequestFromReader(conn)
	if err != nil {
		handlerErr := &HandlerError{
			StatusCode: 400,
			Message:    "Bad Request: " + err.Error(),
		}
		handlerErr.WriteHandlerError(conn)
		return
	}

	// Create a new empty bytes.Buffer for the handler to write to
	buf := bytes.NewBuffer([]byte{})

	// Call the handler function
	handlerErr := s.handler(buf, req)

	// If the handler errors, write the error to the connection
	if handlerErr != nil {
		handlerErr.WriteHandlerError(conn)
		return
	}

	// If the handler succeeds:
	// Get the response body from the handler's buffer
	body := buf.Bytes()

	// Create new default response headers
	headers := response.GetDefaultHeaders(len(body))

	// Write the status line
	response.WriteStatusLine(conn, response.StatusCodeOK)
	// Write the headers
	response.WriteHeaders(conn, headers)
	// write the body
	conn.Write(body)
}
