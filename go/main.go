package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type synchronizeDirectoryTreesArguments struct {
	SourceDirectory      string `json:"sourceDirectory"`
	DestinationDirectory string `json:"destinationDirectory"`
}

func main() {
	if len(os.Args) > 1 {
		jsonString := os.Args[1]
		var arguments synchronizeDirectoryTreesArguments
		if err := json.Unmarshal([]byte(jsonString), &arguments); err != nil {
			fmt.Printf("Error parsing JSON: %s\n", err)
			return
		}
		fmt.Printf("Received message: %s\n", arguments.SourceDirectory)
	} else {
		fmt.Println("No message received.")
	}
}
