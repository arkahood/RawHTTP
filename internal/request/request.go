package request

import (
	"RAWHTTP/internal/headers"
	"errors"
	"fmt"
	"io"
	"strings"
)

const crlf = "\r\n"
const bufferSize int = 8

type RequestState int

const (
	RequestStateInitialized RequestState = iota
	RequestStateDone
	RequestStateParsingHeaders
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	State       RequestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	req := &Request{
		State:   RequestStateInitialized,
		Headers: make(headers.Headers),
	}

	for req.State != RequestStateDone {
		// If the buffer is full, grow it
		if readToIndex == cap(buf) {
			newBuf := make([]byte, cap(buf)*2)
			copy(newBuf, buf[:readToIndex])
			buf = newBuf
		}

		// Read from the io.Reader into the buffer starting at readToIndex
		n, err := reader.Read(buf[readToIndex:])

		if err != nil {
			if errors.Is(err, io.EOF) {
				req.State = RequestStateDone
				break
			}
			return nil, err
		}

		// Update readToIndex with the number of bytes actually read
		readToIndex += n

		// Call parse with the slice of buffer that has data we've read so far
		bytesConsumed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		// Remove the data that was parsed successfully from the buffer
		if bytesConsumed > 0 {
			remainingData := buf[bytesConsumed:readToIndex]
			newBuf := make([]byte, cap(buf))
			copy(newBuf, remainingData)
			buf = newBuf

			// Decrement readToIndex by the number of bytes that were parsed
			readToIndex -= bytesConsumed
		}
	}

	return req, nil
}

func parseRequestLine(data string) (*RequestLine, int, error) {
	// Look for the CRLF that marks the end of the request line
	endIdx := strings.Index(data, crlf)
	if endIdx == -1 {
		// Not enough data yet; no CRLF found
		return nil, 0, nil
	}

	// Extract the request line (without the trailing CRLF)
	reqLine := data[:endIdx]

	parts := strings.Split(reqLine, " ")
	if len(parts) != 3 {
		return nil, endIdx + 2, errors.New("invalid number of parts in request line")
	}
	//  "method" part only contains capital alphabetic characters.
	if strings.ToUpper(parts[0]) != parts[0] {
		return nil, endIdx + 2, errors.New("http method is not capitalized")
	}

	httpVersion := strings.Replace(parts[2], "HTTP/", "", 1)

	if httpVersion != "1.1" {
		return nil, endIdx + 2, errors.New("http/1.1 only supported")
	}

	return &RequestLine{
		Method:        parts[0],
		HttpVersion:   httpVersion,
		RequestTarget: parts[1],
	}, endIdx + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.State != RequestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return totalBytesParsed, err
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	if r.State == RequestStateDone {
		return 0, errors.New("trying to read from done state")
	}
	if r.State == RequestStateInitialized {
		requestLine, numOfByteConsumed, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		// if zero bytes are parsed and no error
		if numOfByteConsumed == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.State = RequestStateParsingHeaders
		return numOfByteConsumed, nil
	}

	if r.State == RequestStateParsingHeaders {
		numberOfBytes, done, err := r.Headers.Parse(data)
		if err != nil {
			fmt.Println("error occured while parsing headers")
			return 0, err
		}
		if done {
			r.State = RequestStateDone
		}
		return numberOfBytes, nil
	}

	return 0, errors.New("unknown request status")
}
