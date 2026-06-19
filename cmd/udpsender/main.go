package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatalf("Failed to resolve udp address: %v\n", err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Failed to dial udp: %v\n", err)
	}
	defer conn.Close()
	buff := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(">")
		line, err := buff.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			fmt.Printf("Error read buff stdin :%v\n", err)
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Printf("Error write to UDP: %v\n", err)
		}
	}
}
