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

func jsonToFunctionCallWithArguments(jsonFunctionCall, jsonArguments string) string {
	var functionCall string
	err := json.Unmarshal([]byte(jsonFunctionCall), &functionCall)
	if err == nil {
		// TODO: enum
		switch functionCall {
		case "synchronizeDirectoryTreesToJSON":
			var arguments SynchronizeDirectoryTreesArguments
			if err = json.Unmarshal([]byte(jsonArguments), &arguments); err == nil {
				return synchronizeDirectoryTreesToJSON(arguments.SourceDirectory, arguments.DestinationDirectory)
			}
		}
	}
	return errorMessageToJSONFunctionResult(err.Error())
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
