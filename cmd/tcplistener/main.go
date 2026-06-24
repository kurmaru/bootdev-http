package main

import (
	"fmt"
	"log"
	"net"

	"github.com/kurmaru/bootdev-http/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Can't listen: %v\n", err)
	}
	defer listener.Close()
	fmt.Println("Listening on port 42069")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Connection failed: %v\n", err)
		}

		fmt.Printf("Connection accepted\n")
		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("Parse from connection failed: %v\n", err)
		}

		fmt.Printf(
			"Request line:\n- Method: %v\n- Target: %v\n- Version: %v\n",
			req.RequestLine.Method,
			req.RequestLine.RequestTarget,
			req.RequestLine.HttpVersion,
		)
		fmt.Println("Headers:")
		for key, val := range req.Headers {
			fmt.Printf("- %v: %v\n", key, val)
		}
		fmt.Printf("Body:\n%s\n", req.Body)
	}
}
