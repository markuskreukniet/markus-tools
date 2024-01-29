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

func stringsToFunctionCallWithArguments(functionCall, jsonArguments string) string {
	var err error
	switch functionCall {
	case FunctionCallSynchronizeDirectoryTreesToJSON:
		var arguments SynchronizeDirectoryTreesArguments
		if err = json.Unmarshal([]byte(jsonArguments), &arguments); err == nil {
			return synchronizeDirectoryTreesToJSON(arguments.SourceDirectory, arguments.DestinationDirectory)
		}
	case FunctionCallGetDuplicateFilesAsNewlineSeparatedStringToJSON:
		var argument GetDuplicateFilesAsNewlineSeparatedStringArgument
		if err = json.Unmarshal([]byte(jsonArguments), &argument); err == nil {
			return getDuplicateFilesAsNewlineSeparatedStringToJSON(argument.UniqueFileSystemNodes)
		}
	}
	errorMessage := "did not receive a correct function call string"
	if err != nil {
		errorMessage = err.Error()
	}
	return errorMessageToJSONFunctionResult(errorMessage)
}

// TODO: check if the logic with starting with and without capitals is correct, for example for the functions and vars
// TODO: if if to if else, where it makes sense. It does not make sense with an 'if err != nil' check
func main() {
	if len(os.Args) > 2 {
		fmt.Print(stringsToFunctionCallWithArguments(os.Args[1], os.Args[2]))
	} else {
		fmt.Print(errorMessageToJSONFunctionResult("os.Args did not receive at least three arguments"))
	}
}
