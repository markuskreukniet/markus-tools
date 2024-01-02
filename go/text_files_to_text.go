package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func isZeroByteFileATextFile(filePath string) (bool, error) {
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
		builder.WriteString(scanner.Text())
	}
	for scanner.Scan() {
		builder.WriteString("\n" + scanner.Text())
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

func textFilesToText(uniqueFileSystemNodes []FileSystemNode) (string, error) {
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
			// TODO: duplicate: if fileDetail.Size > 0 {
			if fileDetail.Size > 0 {
				isTextFile, err := isZeroByteFileATextFile(fileDetail.Path)
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
