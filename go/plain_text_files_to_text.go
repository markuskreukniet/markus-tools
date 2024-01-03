package main

import (
	"bufio"
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

	// Read the first 512 bytes to check for non-text characters. DetectContentType of package 'net/http' works with a similar check.
	bytes := make([]byte, 512)
	_, err = file.Read(bytes)
	if err != nil {
		return false, err
	}
	for _, byte := range bytes {
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
	if scanner.Scan() {
		_, err := builder.WriteString(scanner.Text())
		if err != nil {
			return err
		}
	}
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

func plainTextFilesToText(uniqueFileSystemNodes []FileSystemNode) (string, error) {
	var filePaths []string
	for _, node := range uniqueFileSystemNodes {
		if node.IsDirectory {
			err := walkFileDetails(node.Path, filesWithoutZeroByteFiles, plainTextFiles, func(fileDetail FileDetail) {
				filePaths = append(filePaths, fileDetail.Path)
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
		// TODO: use this [1:] also on other places
		for _, path := range filePaths[1:] {
			err := addLastPathElementAndAllLinesToBuilder("\n"+path, &result)
			if err != nil {
				return "", err
			}
		}
	}
	return result.String(), nil
}
