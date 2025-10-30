package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)

		currentLine := ""

		for {
			read := make([]byte, 8)
			_, err := f.Read(read)

			if err == io.EOF {
				// Send the last line if it has content
				if currentLine != "" {
					lines <- currentLine
				}
				return
			}

			if err != nil {
				fmt.Println("Some Error Happened While Putting in Slice")
				return
			}

			parts := strings.Split(string(read), "\n")

			// Process all parts except the last one, which may be incomplete
			if len(parts)-1 > 0 {
				currentLine += parts[0]
				lines <- currentLine
				currentLine = "" // Reset for next line
			}

			// The last part, which may be incomplete gets added to currentLine
			currentLine += parts[len(parts)-1]
		}
	}()

	return lines
}

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
			lines := getLinesChannel(conn)

			for line := range lines {
				fmt.Println(line)
			}
		}(conn)
	}
}
