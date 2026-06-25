// Package server extend tcp connection, handle incomming request
package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/kurmaru/bootdev-http/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		return nil, err
	}

	s := Server{
		listener: l,
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

	err := response.WriteStatusLine(conn, 200)
	if err != nil {
		fmt.Printf("failed to write status line: %v\n", err)
		return
	}
	err = response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	if err != nil {
		fmt.Printf("failed to write headers: %v\n", err)
	}
}
