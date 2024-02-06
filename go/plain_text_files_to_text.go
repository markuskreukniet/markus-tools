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
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		_, err := utils.WriteNewlineString(builder)
		if err != nil {
			return err
		}
		_, err = builder.WriteString(scanner.Text())
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}

func addLastPathElementAndAllLinesToBuilder(filePath string, builder *strings.Builder) error {
	_, err := builder.WriteString(filepath.Base(filePath))
	if err != nil {
		return err
	}
	return readLinesAddToBuilder(filePath, builder)
}

func plainTextFilesToTextToJSON(uniqueFileSystemNodes []fileSystemNode) string {
	text, err := plainTextFilesToText(uniqueFileSystemNodes)
	if err != nil {
		return errorToJSONFunctionResult(err)
	}
	return resultToJSONFunctionResult(text)
}

// Opening a file two times is not the most efficient, but having a separate open file in isNonZeroByteFileATextFile helps with filtering.
func plainTextFilesToText(uniqueFileSystemNodes []fileSystemNode) (string, error) {
	var filePaths []string
	for _, node := range uniqueFileSystemNodes {
		if node.IsDirectory {
			err := utils.WalkFileDetails(node.Path, utils.FilesWithoutZeroByteFiles, utils.PlainTextFiles, func(detail utils.FileDetail) {
				filePaths = append(filePaths, detail.Path)
			})
			if err != nil {
				return "", err
			}
		} else {
			fileDetail, err := utils.GetFileDetail(node.Path)
			if err != nil {
				return "", err
			}
			if utils.IsFileDetailNonZeroByte(fileDetail) {
				isTextFile, err := utils.IsNonZeroByteFileATextFile(fileDetail.Path)
				if err != nil {
					return "", err
				}
				if isTextFile {
					filePaths = append(filePaths, fileDetail.Path)
				}
			}
		}
	}
	var result strings.Builder
	if len(filePaths) > 0 {
		err := addLastPathElementAndAllLinesToBuilder(filePaths[0], &result)
		if err != nil {
			return "", err
		}
		for i := 1; i < len(filePaths); i++ {
			_, err := result.WriteString("\n\n")
			if err != nil {
				return "", err
			}
			err = addLastPathElementAndAllLinesToBuilder(filePaths[i], &result)
			if err != nil {
				return "", err
			}
		}
	}
	return result.String(), nil
}
