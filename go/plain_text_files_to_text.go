package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func isNonZeroByteFileATextFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read the first 512 or less to check for non-text characters. DetectContentType of package 'net/http' works with a similar check.
	bytes := make([]byte, 512)
	numberOfBytesRead, err := file.Read(bytes)
	if err != nil && err != io.EOF {
		return false, err
	}
	for _, byte := range bytes[:numberOfBytesRead] {
		if !unicode.IsPrint(rune(byte)) && !unicode.IsSpace(rune(byte)) {
			return false, nil
		}
	}
	return true, nil
}

func readLinesAddToBuilder(filePath string, builder *strings.Builder) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		_, err := writeNewlineString(builder)
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

// Opening a file two times is not the most efficient, but having a separate open file in isNonZeroByteFileATextFile helps with filtering.
func plainTextFilesToText(uniqueFileSystemNodes []FileSystemNode) (string, error) {
	var filePaths []string
	for _, node := range uniqueFileSystemNodes {
		if node.IsDirectory {
			err := walkFileDetails(node.Path, filesWithoutZeroByteFiles, plainTextFiles, func(detail fileDetail) {
				filePaths = append(filePaths, detail.Path)
			})
			if err != nil {
				return "", err
			}
		} else {
			fileDetail, err := getFileDetail(node.Path)
			if err != nil {
				return "", err
			}
			if isFileDetailNonZeroByte(fileDetail) {
				isTextFile, err := isNonZeroByteFileATextFile(fileDetail.Path)
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
