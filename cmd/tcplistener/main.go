package main

import (
	"RAWHTTP/internal/request"
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")

	if err != nil {
		fmt.Println("Error Happened", err.Error())
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}

		go func(c net.Conn) {
			req, err := request.RequestFromReader(conn)

			if err != nil {
				fmt.Println("error happened", err.Error())
			}

			fmt.Println("Request line:")
			fmt.Println("Method: ", req.RequestLine.Method)
			fmt.Println("Http Version: ", req.RequestLine.HttpVersion)
			fmt.Println("Target: ", req.RequestLine.RequestTarget)
			fmt.Println("Headers")
			for key, val := range req.Headers {
				fmt.Println(key, ": ", val)
			}
		}(conn)
	}
}
