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
	reason := ""
	switch statusCode {
	case OK:
		reason = "OK"
	case BadRequest:
		reason = "Bad Request"
	case InternalServerError:
		reason = "Internal Server Error"
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %v %v\r\n", statusCode, reason)
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

func WriteBody(w io.Writer, data []byte) error {
	_, err := w.Write(data)
	return err
}
