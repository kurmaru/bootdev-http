// Package response handle sending the HTTP response back to client
package response

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/kurmaru/bootdev-http/internal/headers"
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := fmt.Sprintf("HTTP/1.1 %v ", statusCode)
	switch statusCode {
	case OK:
		statusLine += "OK"
	case BadRequest:
		statusLine += "Bad Request"
	case InternalServerError:
		statusLine += "Internal Server Error"
	}
	statusLine += "\r\n"
	_, err := w.Write([]byte(statusLine))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")
	h.Set("Content-Length", strconv.Itoa(contentLen))
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	var str strings.Builder

	for key, val := range headers {
		fmt.Fprintf(&str, "%v: %v\r\n", key, val)
	}
	str.Write([]byte("\r\n"))

	_, err := w.Write([]byte(str.String()))
	if err != nil {
		return err
	}

	return nil
}
