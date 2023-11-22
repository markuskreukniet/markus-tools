package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	FunctionCallSynchronizeDirectoryTreesToJSON                 string = "synchronizeDirectoryTreesToJSON"
	FunctionCallGetDuplicateFilesAsNewlineSeparatedStringToJSON string = "getDuplicateFilesAsNewlineSeparatedStringToJSON"
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
		switch functionCall {
		case FunctionCallSynchronizeDirectoryTreesToJSON:
			var arguments SynchronizeDirectoryTreesArguments
			if err = json.Unmarshal([]byte(jsonArguments), &arguments); err == nil {
				return synchronizeDirectoryTreesToJSON(arguments.SourceDirectory, arguments.DestinationDirectory)
			}
		}
	}
	return errorMessageToJSONFunctionResult(err.Error())
}

func main() {
	if len(os.Args) > 2 {
		fmt.Print(jsonToFunctionCallWithArguments(os.Args[1], os.Args[2]))
	} else {
		fmt.Print(errorMessageToJSONFunctionResult("os.Args did not receive at least three arguments"))
	}
}
