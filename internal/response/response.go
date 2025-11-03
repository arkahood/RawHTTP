package response

import (
	"RAWHTTP/internal/headers"
	"errors"
	"io"
	"strconv"
)

type StatusCode int

const (
	StatusCodeOK StatusCode = iota
	StatusCodeClientErrors
	StatusCodeInternalServerError
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var err error
	switch statusCode {
	case StatusCodeOK:
		_, err = w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case StatusCodeClientErrors:
		_, err = w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case StatusCodeInternalServerError:
		_, err = w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		return errors.New("unsupported response type")
	}
	return err
}
func GetDefaultHeaders(contentLen int) headers.Headers {
	return headers.Headers{
		"content-length": strconv.Itoa(contentLen),
		"content-type":   "text/plain",
		"connection":     "close",
	}
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	var headerRes []byte
	for key, val := range headers {
		headerRes = append(headerRes, []byte(key+":"+val+"\r\n")...)
	}
	headerRes = append(headerRes, []byte("\r\n")...)
	_, err := w.Write(headerRes)
	return err
}
