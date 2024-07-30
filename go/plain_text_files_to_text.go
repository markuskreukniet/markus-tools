package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

func readLinesAddToBuilder(filePath string, builder *strings.Builder) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// TODO: os.ReadFile is better?
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

func addLastPathElementAndAllLinesToBuilder(filePath string, builder *strings.Builder) error {
	if _, err := builder.WriteString(filepath.Base(filePath)); err != nil {
		return err
	}
	return readLinesAddToBuilder(filePath, builder)
}

func plainTextFilesToTextToJSON(uniqueFileSystemNodes []utils.FileSystemNode) string {
	text, err := plainTextFilesToText(uniqueFileSystemNodes)
	if err != nil {
		return errorToJSONFunctionResult(err)
	}
	return resultToJSONFunctionResult(text)
}

// Opening a file two times is not the most efficient, but having a separate open file in isNonZeroByteFileATextFile helps with filtering.
func plainTextFilesToText(uniqueFileSystemNodes []utils.FileSystemNode) (string, error) {
	var filePaths []string

	for _, node := range uniqueFileSystemNodes {
		if node.IsDirectory {
			handler := func(file utils.FileSystemFile) error {
				filePaths = append(filePaths, file.Path)
				return nil
			}

			if err := utils.WalkFilterAndHandleFileSystemFile(node.Path, utils.NonZeroByteFiles, utils.PlainTextFiles, handler); err != nil {
				return "", err
			}
		} else {
			file, err := utils.ToFileSystemFile(node.Path)
			if err != nil {
				return "", err
			}
			if file.FileMetadata.Size > 0 {
				isTextFile, err := utils.IsNonZeroByteFileATextFile(file.Path)
				if err != nil {
					return "", err
				}
				if isTextFile {
					filePaths = append(filePaths, file.Path)
				}
			}
		}
	}

	var result strings.Builder
	length := len(filePaths)

	if length > 0 {
		err := addLastPathElementAndAllLinesToBuilder(filePaths[0], &result)
		if err != nil {
			return "", err
		}
		for i := 1; i < length; i++ {
			if _, err := utils.WriteTwoNewlineStrings(&result); err != nil {
				return "", err
			}
			if err = addLastPathElementAndAllLinesToBuilder(filePaths[i], &result); err != nil {
				return "", err
			}
		}
	}

	return result.String(), nil
}
