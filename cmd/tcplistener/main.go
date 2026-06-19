package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
		linesCh := getLinesChannel(conn)
		for line := range linesCh {
			fmt.Printf("%v\n", line)
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)
		curStr := ""
		for {
			buff := [8]byte{}
			count, err := f.Read(buff[:])
			if err != nil {
				log.Fatalf("failed to read file :%v\n", err)
			}

			parts := strings.Split(string(buff[:]), "\n")
			for i, line := range parts {
				curStr += line
				if len(parts) > 1 && i < len(parts)-1 {
					ch <- curStr
					curStr = ""
				}
			}

			if count < 8 {
				break
			}
		}
	}()

	return ch
}
