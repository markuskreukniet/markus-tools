package main

import (
	"strings"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

func getDuplicateFilesAsNewlineSeparatedStringToJSON(uniqueFileSystemNodes []utils.FileSystemNode) string {
	newlineSeparatedString, err := getDuplicateFilesAsNewlineSeparatedString(uniqueFileSystemNodes)
	if err != nil {
		return errorToJSONFunctionResult(err)
	}

	return resultToJSONFunctionResult(newlineSeparatedString)
}

func getDuplicateFilesAsNewlineSeparatedString(uniqueFileSystemNodes []utils.FileSystemNode) (string, error) {
	var result strings.Builder

	// var files []utils.FileInfo

	// // is the error return needed?
	// handler := func(file utils.FileInfo) error {
	// 	files = append(files, utils.MinimalFileInfo{
	// 		Size:         file.GetSize(),
	// 		AbsolutePath: file.GetAbsolutePath(),
	// 	})
	// 	return nil
	// }

	// for _, node := range uniqueFileSystemNodes {
	// 	utils.WalkFilterAndHandleFileInfo(node, utils.NonZeroByteFiles, utils.AllFiles, handler)
	// }

	// TODO: should be FileMetadata
	var files []utils.FileSystemFile
	if err := utils.AppendNonZeroByteFiles(uniqueFileSystemNodes, &files); err != nil {
		return "", err
	}

	groups, err := utils.CreateFileSystemFileByHashGroups(files, true)
	if err != nil {
		return "", err
	}

	for i, group := range groups {
		if i != 0 {
			if _, err := utils.WriteTwoNewlineStrings(&result); err != nil {
				return "", err
			}
		}
		for j, file := range group {
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
