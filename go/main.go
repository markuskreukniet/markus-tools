package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// TODO: there are duplicate or useless things, such as statements, strings, and structs, probably in tests

const (
	functionCallSynchronizeDirectoryTreesToJSON                 string = "synchronizeDirectoryTreesToJSON"
	functionCallGetDuplicateFilesAsNewlineSeparatedStringToJSON string = "getDuplicateFilesAsNewlineSeparatedStringToJSON"
)

type getDuplicateFilesAsNewlineSeparatedStringArgument struct {
	UniqueFileSystemNodes []fileSystemNode `json:"uniqueFileSystemNodes"`
}

type synchronizeDirectoryTreesArguments struct {
	SourceDirectory      string `json:"sourceDirectory"`
	DestinationDirectory string `json:"destinationDirectory"`
}

func stringsToFunctionCallWithArguments(functionCall, jsonArguments string) string {
	var err error
	switch functionCall {
	case functionCallSynchronizeDirectoryTreesToJSON:
		var arguments synchronizeDirectoryTreesArguments
		if err = json.Unmarshal([]byte(jsonArguments), &arguments); err == nil {
			return synchronizeDirectoryTreesToJSON(arguments.SourceDirectory, arguments.DestinationDirectory)
		}
	case functionCallGetDuplicateFilesAsNewlineSeparatedStringToJSON:
		var argument getDuplicateFilesAsNewlineSeparatedStringArgument
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

func main() {
	if len(os.Args) > 2 {
		fmt.Print(stringsToFunctionCallWithArguments(os.Args[1], os.Args[2]))
	} else {
		fmt.Print(errorMessageToJSONFunctionResult("os.Args did not receive at least three arguments"))
	}
}
