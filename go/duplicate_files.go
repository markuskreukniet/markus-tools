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
	if _, err := result.WriteString(duplicateFiles[0].path); err != nil {
		return "", err
	}
	for i := 1; i < len(duplicateFiles); i++ {
		if _, err := utils.WriteNewlineString(&result); err != nil {
			return "", err
		}
		if duplicateFiles[i].hash != duplicateFiles[i-1].hash {
			if _, err := utils.WriteNewlineString(&result); err != nil {
				return "", err
			}
		}
		if _, err := result.WriteString(duplicateFiles[i].path); err != nil {
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
	if err := utils.AppendFileDetails(
		func(detail utils.FileDetail) {
			fileIdentifiers = append(fileIdentifiers, fileIdentifier{
				path: detail.Path,
				size: detail.Size,
				hash: "",
			})
		}, uniqueFileSystemNodes, utils.FilesWithoutZeroByteFiles); err != nil {
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
				if fileIdentifiers[previousIndex].hash, err = getFileHash(fileIdentifiers[previousIndex].path); err != nil {
					return "", err
				}
			}
			if fileIdentifiers[i].hash, err = getFileHash(fileIdentifiers[i].path); err != nil {
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
