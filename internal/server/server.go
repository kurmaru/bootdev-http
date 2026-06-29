// Package server extend tcp connection, handle incomming request
package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/kurmaru/bootdev-http/internal/request"
	"github.com/kurmaru/bootdev-http/internal/response"
)

type Server struct {
	handler  Handler
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		return nil, err
	}

	s := Server{
		listener: l,
		handler:  handler,
	}

	go s.listen()

	return &s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			fmt.Printf("failed to accept connection: %v\n", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	res := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("Failed to parse: %v", err)

		body := fmt.Sprintf("Parse from connection failed: %v\n", err)
		res.WriteStatusLine(response.BadRequest)
		res.WriteHeaders(response.GetDefaultHeaders(len(body)))
		res.WriteBody([]byte(body))
		return
	}

	s.handler(*res, *req)
}
