package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

type fileIdentifier struct {
	path string
	size int64
	hash string
}

type duplicateFile struct {
	path string
	hash string
}

// TODO: is this struct needed?
type appendFileIdentifier struct {
	fileIdentifiers *[]fileIdentifier
}

func (identifier *appendFileIdentifier) Append(detail utils.FileDetail) {
	*identifier.fileIdentifiers = append(*identifier.fileIdentifiers, fileIdentifier{
		path: detail.Path,
		size: detail.Size,
		hash: "",
	})
}

func appendDuplicateFile(duplicateFiles *[]duplicateFile, fileIdentifier fileIdentifier) {
	*duplicateFiles = append(*duplicateFiles, duplicateFile{
		path: fileIdentifier.path,
		hash: fileIdentifier.hash,
	})
}

func getFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// sha256 is generally faster and more secure than SHA1
	// SHA1 is generally faster and more secure than MD5
	hashGenerator := sha256.New()
	if _, err := io.Copy(hashGenerator, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hashGenerator.Sum(nil)), nil
}

func duplicateFilesToNewlineSeparatedString(duplicateFiles []duplicateFile) (string, error) {
	if len(duplicateFiles) == 0 {
		return "", nil
	}
	var result strings.Builder
	_, err := result.WriteString(duplicateFiles[0].path)
	if err != nil {
		return "", err
	}
	for i := 1; i < len(duplicateFiles); i++ {
		_, err = utils.WriteNewlineString(&result)
		if err != nil {
			return "", err
		}
		if duplicateFiles[i].hash != duplicateFiles[i-1].hash {
			_, err = utils.WriteNewlineString(&result)
			if err != nil {
				return "", err
			}
		}
		_, err = result.WriteString(duplicateFiles[i].path)
		if err != nil {
			return "", err
		}
	}
	return result.String(), nil
}

func getDuplicateFilesAsNewlineSeparatedStringToJSON(uniqueFileSystemNodes []utils.FileSystemNode) string {
	newlineSeparatedString, err := getDuplicateFilesAsNewlineSeparatedString(uniqueFileSystemNodes)
	if err != nil {
		return errorToJSONFunctionResult(err)
	}
	return resultToJSONFunctionResult(newlineSeparatedString)
}

func getDuplicateFilesAsNewlineSeparatedString(uniqueFileSystemNodes []utils.FileSystemNode) (string, error) {
	var fileIdentifiers []fileIdentifier
	err := utils.AppendFileDetails(&appendFileIdentifier{fileIdentifiers: &fileIdentifiers}, uniqueFileSystemNodes)
	if err != nil {
		return "", err
	}
	sort.Slice(fileIdentifiers, func(i, j int) bool {
		return fileIdentifiers[i].size < fileIdentifiers[j].size
	})
	// TODO: We can make this function efficient by making from here the newlineSeparatedString, instead of making first the duplicateFiles slice
	var duplicateFiles []duplicateFile
	var lastAppendedIndex = -1
	for i := 1; i < len(fileIdentifiers); i++ {
		previousIndex := i - 1
		if fileIdentifiers[previousIndex].size == fileIdentifiers[i].size {
			var err error
			if fileIdentifiers[previousIndex].hash == "" {
				fileIdentifiers[previousIndex].hash, err = getFileHash(fileIdentifiers[previousIndex].path)
				if err != nil {
					return "", err
				}
			}
			fileIdentifiers[i].hash, err = getFileHash(fileIdentifiers[i].path)
			if err != nil {
				return "", err
			}
			if fileIdentifiers[previousIndex].hash == fileIdentifiers[i].hash {
				if lastAppendedIndex != previousIndex {
					appendDuplicateFile(&duplicateFiles, fileIdentifiers[previousIndex])
				}
				appendDuplicateFile(&duplicateFiles, fileIdentifiers[i])
				lastAppendedIndex = i
			}
		}
	}
	newlineSeparatedString, err := duplicateFilesToNewlineSeparatedString(duplicateFiles)
	if err != nil {
		return "", err
	}
	return newlineSeparatedString, nil
}
