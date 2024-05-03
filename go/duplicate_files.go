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

type duplicateFileGroup struct {
	hash      string
	filePaths []string
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

	// create duplicate file groups
	var duplicateFileGroups []duplicateFileGroup
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
			foundGroup := false
			for j, group := range duplicateFileGroups {
				if fileIdentifiers[i].hash == group.hash {
					foundGroup = true
					duplicateFileGroups[j].filePaths = append(duplicateFileGroups[j].filePaths, fileIdentifiers[i].path)
					break
				}
			}
			if !foundGroup {
				for j := 0; j <= previousIndex; j++ {
					if fileIdentifiers[i].hash == fileIdentifiers[j].hash {
						duplicateFileGroups = append(duplicateFileGroups, duplicateFileGroup{
							hash:      fileIdentifiers[i].hash,
							filePaths: []string{fileIdentifiers[j].path, fileIdentifiers[i].path},
						})
						break
					}
				}
			}
		}
	}

	// create and return the result string
	var result strings.Builder
	for i, group := range duplicateFileGroups {
		if i != 0 {
			if _, err := utils.WriteTwoNewlineStrings(&result); err != nil {
				return "", err
			}
		}
		for j, path := range group.filePaths {
			if j != 0 {
				if _, err := utils.WriteNewlineString(&result); err != nil {
					return "", err
				}
			}
			if _, err := result.WriteString(path); err != nil {
				return "", err
			}
		}
	}
	return result.String(), nil
}
