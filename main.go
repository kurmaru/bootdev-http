package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatalf("failed to open file :%v\n", err)
	}

	curStr := ""

	for {
		buff := [8]byte{}
		count, err := file.Read(buff[:])
		if err != nil {
			log.Fatalf("failed to read file :%v\n", err)
		}

		parts := strings.Split(string(buff[:]), "\n")

		for i, line := range parts {
			curStr += line
			if len(parts) > 1 && i < len(parts)-1 {
				fmt.Printf("read: %v\n", curStr)
				curStr = ""
			}
		}

		if count < 8 {
			break
		}
	}
}
