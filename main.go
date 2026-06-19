package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatalf("failed to open file :%v\n", err)
	}

	readCh := getLinesChannel(file)
	for line := range readCh {
		fmt.Printf("read: %v\n", line)
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
