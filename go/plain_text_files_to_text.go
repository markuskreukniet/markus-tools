package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

func addReadLines(filePath string, builder *strings.Builder) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if _, err := utils.WriteNewlineString(builder); err != nil {
			return err
		}
		if _, err = builder.WriteString(scanner.Text()); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func addFilePathBaseAndAllLinesToBuilder(filePath string, builder *strings.Builder) error {
	if _, err := builder.WriteString(filepath.Base(filePath)); err != nil {
		return err
	}
	return addReadLines(filePath, builder)
}

func plainTextFilesToTextToJSON(uniqueFileSystemNodes []utils.FileSystemNode) string {
	text, err := plainTextFilesToText(uniqueFileSystemNodes)
	if err != nil {
		return errorToJSONFunctionResult(err)
	}
	return resultToJSONFunctionResult(text)
}

// Opening a file two times is not the most efficient, but having a separate open file in isTextFile helps with filtering.
func plainTextFilesToText(uniqueFileSystemNodes []utils.FileSystemNode) (string, error) {
	var filePaths []string

	for _, node := range uniqueFileSystemNodes {
		if node.IsDirectory {
			handler := func(file utils.FileSystemFile) error {
				filePaths = append(filePaths, file.FileMetadata.Path)
				return nil
			}

			if err := utils.WalkFilterAndHandleFileSystemFile(node.Path, utils.NonZeroByteFiles, utils.TextFiles, handler); err != nil {
				return "", err
			}
		} else {
			file, err := utils.ToFileSystemFile(node.Path)
			if err != nil {
				return "", err
			}
			if file.FileMetadata.Size > 0 {
				isTextFile, err := utils.IsTextFile(file.FileMetadata.Path)
				if err != nil {
					return "", err
				}
				if isTextFile {
					filePaths = append(filePaths, file.FileMetadata.Path)
				}
			}
		}
	}

	var result strings.Builder
	length := len(filePaths)

	if length > 0 {
		err := addFilePathBaseAndAllLinesToBuilder(filePaths[0], &result)
		if err != nil {
			return "", err
		}
		for _, path := range filePaths[1:] {
			if _, err := utils.WriteTwoNewlineStrings(&result); err != nil {
				return "", err
			}
			if err = addFilePathBaseAndAllLinesToBuilder(path, &result); err != nil {
				return "", err
			}
		}
	}

	return result.String(), nil
}
