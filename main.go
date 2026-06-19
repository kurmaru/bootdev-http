package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatalf("failed to open file :%v\n", err)
	}

	for {
		buff := [8]byte{}
		count, err := file.Read(buff[:])
		if err != nil {
			log.Fatalf("failed to read file :%v\n", err)
		}
		fmt.Printf("read: %v\n", string(buff[:]))
		if count < 8 {
			break
		}
	}
}
