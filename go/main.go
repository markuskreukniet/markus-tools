package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type MyData struct {
	SourceDirectoryFilePath      string `json:"sourceDirectoryFilePath"`
	DestinationDirectoryFilePath string `json:"destinationDirectoryFilePath"`
}

func main() {
	if len(os.Args) > 1 {
		jsonData := os.Args[1]

		// Parse the JSON data into the MyData struct
		var data MyData
		if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
			fmt.Printf("Error parsing JSON: %s\n", err)
			return
		}

		fmt.Printf("Received message: %s\n", data.SourceDirectoryFilePath)
	} else {
		fmt.Println("No message received.")
	}
}
