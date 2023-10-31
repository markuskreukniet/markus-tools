package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		message := os.Args[1]
		fmt.Printf("Received message: %s\n", message)
	} else {
		fmt.Println("No message received.")
	}
}
