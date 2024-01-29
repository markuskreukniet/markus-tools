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

type FileSystemNode struct {
	Path        string
	IsDirectory bool
}

type FileIdentifier struct {
	Path string
	Size int64
	Hash string
}

type DuplicateFile struct {
	Path string
	Hash string
}

func appendFileIdentifier(fileIdentifiers *[]FileIdentifier, detail utils.FileDetail) {
	*fileIdentifiers = append(*fileIdentifiers, FileIdentifier{
		Path: detail.Path,
		Size: detail.Size,
		Hash: "",
	})
}

func appendDuplicateFile(duplicateFiles *[]DuplicateFile, fileIdentifier FileIdentifier) {
	*duplicateFiles = append(*duplicateFiles, DuplicateFile{
		Path: fileIdentifier.Path,
		Hash: fileIdentifier.Hash,
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

func duplicateFilesToNewlineSeparatedString(duplicateFiles []DuplicateFile) (string, error) {
	if len(duplicateFiles) == 0 {
		return "", nil
	}
	var result strings.Builder
	_, err := result.WriteString(duplicateFiles[0].Path)
	if err != nil {
		return "", err
	}
	for i := 1; i < len(duplicateFiles); i++ {
		_, err = utils.WriteNewlineString(&result)
		if err != nil {
			return "", err
		}
		if duplicateFiles[i].Hash != duplicateFiles[i-1].Hash {
			_, err = utils.WriteNewlineString(&result)
			if err != nil {
				return "", err
			}
		}
		_, err = result.WriteString(duplicateFiles[i].Path)
		if err != nil {
			return "", err
		}
	}
	return result.String(), nil
}

func getDuplicateFilesAsNewlineSeparatedStringToJSON(uniqueFileSystemNodes []FileSystemNode) string {
	newlineSeparatedString, err := getDuplicateFilesAsNewlineSeparatedString(uniqueFileSystemNodes)
	if err != nil {
		return errorMessageToJSONFunctionResult(err.Error())
	}
	return resultToJSONFunctionResult(newlineSeparatedString)
}

func getDuplicateFilesAsNewlineSeparatedString(uniqueFileSystemNodes []FileSystemNode) (string, error) {
	var fileIdentifiers []FileIdentifier
	for _, node := range uniqueFileSystemNodes {
		if node.IsDirectory {
			err := utils.WalkFileDetails(node.Path, utils.FilesWithoutZeroByteFiles, utils.AllFiles, func(detail utils.FileDetail) {
				appendFileIdentifier(&fileIdentifiers, detail)
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
				appendFileIdentifier(&fileIdentifiers, fileDetail)
			}
		}
	}
	sort.Slice(fileIdentifiers, func(i, j int) bool {
		return fileIdentifiers[i].Size < fileIdentifiers[j].Size
	})
	// TODO: We can make this function efficient by making from here the newlineSeparatedString, instead of making first the duplicateFiles slice
	var duplicateFiles []DuplicateFile
	var lastAppendedIndex = -1
	for i := 1; i < len(fileIdentifiers); i++ {
		previousIndex := i - 1
		if fileIdentifiers[previousIndex].Size == fileIdentifiers[i].Size {
			var err error
			if fileIdentifiers[previousIndex].Hash == "" {
				fileIdentifiers[previousIndex].Hash, err = getFileHash(fileIdentifiers[previousIndex].Path)
				if err != nil {
					return "", err
				}
			}
			fileIdentifiers[i].Hash, err = getFileHash(fileIdentifiers[i].Path)
			if err != nil {
				return "", err
			}
			if fileIdentifiers[previousIndex].Hash == fileIdentifiers[i].Hash {
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
