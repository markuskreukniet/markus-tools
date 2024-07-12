package main

import (
	"strings"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

// duplicateFileGroups
type duplicateFileGroup struct {
	identifier string
	filePaths  []string
}

type duplicateFileGroups []duplicateFileGroup

func (groups duplicateFileGroups) AppendByIdentifier(identifier, filePath string) bool {
	for i, group := range groups {
		if identifier == group.identifier {
			groups[i].filePaths = append(groups[i].filePaths, filePath)
			return true
		}
	}
	return false
}

func getDuplicateFilesAsNewlineSeparatedStringToJSON(uniqueFileSystemNodes []utils.FileSystemNode) string {
	newlineSeparatedString, err := getDuplicateFilesAsNewlineSeparatedString(uniqueFileSystemNodes)
	if err != nil {
		return errorToJSONFunctionResult(err)
	}

	return resultToJSONFunctionResult(newlineSeparatedString)
}

func getDuplicateFilesAsNewlineSeparatedString(nodes []utils.FileSystemNode) (string, error) {
	var result strings.Builder
	var files []utils.FileData
	handler := func(metadata utils.FileMetadata) {
		files = append(files, utils.CreateFileData("", metadata))
	}

	if err := utils.FilterAndHandleAllNodesFileMetadata(nodes, utils.FilesWithoutZeroByteFiles, handler); err != nil {
		return "", err
	}

	groups, err := utils.CreateFileHashGroups(files, true)
	if err != nil {
		return "", err
	}

	for i, group := range groups {
		if i != 0 {
			if _, err := utils.WriteTwoNewlineStrings(&result); err != nil {
				return "", err
			}
		}
		for j, file := range group.FilesData {
			if j != 0 {
				if _, err := utils.WriteNewlineString(&result); err != nil {
					return "", err
				}
			}
			if _, err := result.WriteString(file.FileMetadata.Path); err != nil {
				return "", err
			}
		}
	}

	return result.String(), nil
}
