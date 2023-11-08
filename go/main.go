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
	var result = jsonMarshalWithFallbackJSONError("os.Args did not receive at least two arguments")
	if len(os.Args) > 1 {
		var arguments synchronizeDirectoryTreesArguments
		if err := json.Unmarshal([]byte(os.Args[1]), &arguments); err != nil {
			result = jsonMarshalWithFallbackJSONError(err.Error())
		} else {
			result = arguments.SourceDirectory + " test"
			// result = internal.SynchronizeDirectoryTreesToJSON(arguments.SourceDirectory, arguments.DestinationDirectory)
		}
	}
	fmt.Print(result)
}
