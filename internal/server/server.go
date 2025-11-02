package server

import (
	"errors"
	"fmt"
	"net"
	"sync/atomic"
)

type Server struct {
	listener       net.Listener
	serverIsClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	lstnr, err := net.Listen("tcp", fmt.Sprint(":", port))
	if err != nil {
		return nil, errors.New("tcp error happened")
	}
	server := Server{
		listener: lstnr,
	}

	go func(s *Server) {
		s.listen()
	}(&server)

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

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 12\r\n" +
		"Connection: close\r\n" +
		"\r\n" +
		"Hello World!"

	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Write error:", err)
	}
}
