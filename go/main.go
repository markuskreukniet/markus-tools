package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	FunctionCallSynchronizeDirectoryTreesToJSON string = "synchronizeDirectoryTreesToJSON"
	// FunctionCallGetDuplicateFilesAsNewlineSeparatedStringToJSON string = "getDuplicateFilesAsNewlineSeparatedStringToJSON"
)

// type GetDuplicateFilesAsNewlineSeparatedStringArgument struct {
// 	UniqueFileSystemNodes []FileSystemNode `json:"uniqueFileSystemNodes"`
// }

type SynchronizeDirectoryTreesArguments struct {
	SourceDirectory      string `json:"sourceDirectory"`
	DestinationDirectory string `json:"destinationDirectory"`
}

func stringsToFunctionCallWithArguments(functionCall, jsonArguments string) string {
	var err error
	switch functionCall {
	case FunctionCallSynchronizeDirectoryTreesToJSON:
		var arguments SynchronizeDirectoryTreesArguments
		if err = json.Unmarshal([]byte(jsonArguments), &arguments); err == nil {
			return synchronizeDirectoryTreesToJSON(arguments.SourceDirectory, arguments.DestinationDirectory)
		}
	}
	errorMessage := "did not receive a correct function call string"
	if err != nil {
		errorMessage = err.Error()
	}
	return errorMessageToJSONFunctionResult(errorMessage)
}

func main() {
	if len(os.Args) > 2 {
		fmt.Print(stringsToFunctionCallWithArguments(os.Args[1], os.Args[2]))
	} else {
		fmt.Print(errorMessageToJSONFunctionResult("os.Args did not receive at least three arguments"))
	}
}
