package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

// TODO: rename functionCall to something shorter, to func?
const (
	functionCallSynchronizeDirectoryTreesToJSON                 string = "synchronizeDirectoryTreesToJSON"
	functionCallGetDuplicateFilesAsNewlineSeparatedStringToJSON string = "getDuplicateFilesAsNewlineSeparatedStringToJSON"
	functionCallPlainTextFilesToTextToJSON                      string = "plainTextFilesToTextToJSON"
	functionCallFilesToDateRangeDirectoryToJSON                 string = "filesToDateRangeDirectoryToJSON"
	functionCallReferencesByUrlsToJSON                          string = "referencesByUrlsToJSON"
)

type uniqueFileSystemNodes struct {
	UniqueFileSystemNodes []utils.FileSystemNode `json:"uniqueFileSystemNodes"`
}

type synchronizeDirectoryTreesArguments struct {
	SourceDirectoryFilePath      string `json:"sourceDirectoryFilePath"`
	DestinationDirectoryFilePath string `json:"destinationDirectoryFilePath"`
}

type filesToDateRangeDirectoryArguments struct {
	UniqueFileSystemNodes        []utils.FileSystemNode `json:"uniqueFileSystemNodes"` // TODO: use uniqueFileSystemNodes struct?
	DestinationDirectoryFilePath string                 `json:"destinationDirectoryFilePath"`
}

func toFunctionCall(functionCall, jsonArguments string) string {
	var err error
	switch functionCall {
	case functionCallSynchronizeDirectoryTreesToJSON:
		var arguments synchronizeDirectoryTreesArguments
		if err = json.Unmarshal([]byte(jsonArguments), &arguments); err == nil {
			return synchronizeDirectoryTreesToJSON(
				arguments.SourceDirectoryFilePath, arguments.DestinationDirectoryFilePath,
			)
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
	case functionCallFilesToDateRangeDirectoryToJSON:
		var arguments filesToDateRangeDirectoryArguments
		if err = json.Unmarshal([]byte(jsonArguments), &arguments); err == nil {
			return filesToDateRangeDirectoryToJSON(
				arguments.UniqueFileSystemNodes, arguments.DestinationDirectoryFilePath,
			)
		}
	case functionCallReferencesByUrlsToJSON:
		var argument string
		if err = json.Unmarshal([]byte(jsonArguments), &argument); err == nil {
			return referencesByUrlsToJSON(argument)
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
