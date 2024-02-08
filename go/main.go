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
	functionCallPlainTextFilesToTextToJSON                      string = "functionCallPlainTextFilesToTextToJSON"
)

type uniqueFileSystemNodes struct {
	UniqueFileSystemNodes []fileSystemNode `json:"uniqueFileSystemNodes"`
}

type synchronizeDirectoryTreesArguments struct {
	SourceDirectoryFilePath      string `json:"sourceDirectoryFilePath"`
	DestinationDirectoryFilePath string `json:"destinationDirectoryFilePath"`
}

func toFunctionCall(functionCall, jsonArguments string) string {
	var err error
	switch functionCall {
	case functionCallSynchronizeDirectoryTreesToJSON:
		var arguments synchronizeDirectoryTreesArguments
		if err = json.Unmarshal([]byte(jsonArguments), &arguments); err == nil {
			// TODO: is SourceDirectoryFilePath good naming? check jsx, js, and go files.
			return synchronizeDirectoryTreesToJSON(arguments.SourceDirectoryFilePath, arguments.DestinationDirectoryFilePath)
		}
	case functionCallGetDuplicateFilesAsNewlineSeparatedStringToJSON:
		var argument uniqueFileSystemNodes
		if err = json.Unmarshal([]byte(jsonArguments), &argument); err == nil {
			return getDuplicateFilesAsNewlineSeparatedStringToJSON(argument.UniqueFileSystemNodes)
		}
	case functionCallPlainTextFilesToTextToJSON:
		var argument uniqueFileSystemNodes
		if err = json.Unmarshal([]byte(jsonArguments), &argument); err == nil {
			return plainTextFilesToTextToJSON(argument.UniqueFileSystemNodes)
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
		fmt.Print(toFunctionCall(os.Args[1], os.Args[2]))
	} else {
		fmt.Print(errorMessageToJSONFunctionResult("os.Args did not receive at least three arguments"))
	}
}
