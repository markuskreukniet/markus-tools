package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GetDuplicateFilesAsNewlineSeparatedStringArgument struct {
	UniqueFileSystemNodes []FileSystemNode `json:"uniqueFileSystemNodes"`
}

type SynchronizeDirectoryTreesArguments struct {
	SourceDirectory      string `json:"sourceDirectory"`
	DestinationDirectory string `json:"destinationDirectory"`
}

func main() {
	var result string
	if len(os.Args) > 2 {
		var functionCall string
		if err := json.Unmarshal([]byte(os.Args[1]), &functionCall); err != nil {
			result = errorMessageToJSONFunctionResult(err.Error())
		} else {
			// TODO: enum
			switch functionCall {
			case "SynchronizeDirectoryTrees":
				var arguments SynchronizeDirectoryTreesArguments
				if err = json.Unmarshal([]byte(os.Args[2]), &arguments); err != nil {
					result = errorMessageToJSONFunctionResult(err.Error())
				} else {
					result = synchronizeDirectoryTreesToJSON(arguments.SourceDirectory, arguments.DestinationDirectory)
				}
			}
		}
	} else {
		result = errorMessageToJSONFunctionResult("os.Args did not receive at least three arguments")
	}
	fmt.Print(result)
}
