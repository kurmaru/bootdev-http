// Package server extend tcp connection, handle incomming request
package server

import (
	"bytes"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/kurmaru/bootdev-http/internal/request"
	"github.com/kurmaru/bootdev-http/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int, h Handler) (*Server, error) {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		return nil, err
	}

	s := Server{
		listener: l,
	}

	go s.listen(h)

	return &s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen(h Handler) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			fmt.Printf("failed to accept connection: %v\n", err)
			continue
		}

		go s.handle(conn, h)
	}
}

func (s *Server) handle(conn net.Conn, h Handler) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("Parse from connection failed: %v\n", err)
		return
	}

	buf := new(bytes.Buffer)
	handlerErr := h(buf, *req)
	if handlerErr != nil {
		if err := WriteError(conn, handlerErr); err != nil {
			fmt.Printf("Write to connection failed: %v\n", err)
			return
		}
		return
	}

	if err := response.WriteStatusLine(conn, response.OK); err != nil {
		fmt.Printf("Write status line to connection failed: %v\n", err)
		return
	}

	if err := response.WriteHeaders(conn, response.GetDefaultHeaders(buf.Len())); err != nil {
		fmt.Printf("Write status line to connection failed: %v\n", err)
		return
	}

	if err := response.WriteBody(conn, buf.Bytes()); err != nil {
		fmt.Printf("Write body to connection failed: %v\n", err)
		return
	}
}
