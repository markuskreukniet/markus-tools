package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"sort"
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

func appendFileIdentifier(fileIdentifiers *[]FileIdentifier, fileDetail FileDetail) []FileIdentifier {
	return append(*fileIdentifiers, FileIdentifier{
		Path: fileDetail.Path,
		Size: fileDetail.Size,
		Hash: "",
	})
}

func getFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hashGenerator := sha256.New()
	if _, err := io.Copy(hashGenerator, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hashGenerator.Sum(nil)), nil
}

func duplicateFilesString(uniqueFileSystemNodes []FileSystemNode) (string, error) {
	var fileIdentifiers []FileIdentifier
	for _, value := range uniqueFileSystemNodes {
		if !value.IsDirectory {
			fileDetail, err := getFileDetail(value.Path)
			if err != nil {
				return "", err
			}
			if fileDetail.Size > 0 {
				appendFileIdentifier(&fileIdentifiers, fileDetail)
			}
		} else {
			err := walkFileDetails(value.Path, filesAndDirectoriesWithoutZeroByteFiles, func(fileDetail FileDetail) {
				appendFileIdentifier(&fileIdentifiers, fileDetail)
			})
			if err != nil {
				return "", err
			}
		}
	}
	sort.Slice(fileIdentifiers, func(i, j int) bool {
		return fileIdentifiers[i].Size < fileIdentifiers[j].Size
	})
	var duplicates []FileIdentifier
	var lastAppendedIndex = -1
	// for i := 1; i < len(fileIdentifiers); i++ {
	// 	previousFileIdentifier := fileIdentifiers[i-1] // not needed
	// 	currentFileIdentifier := fileIdentifiers[i]    // not needed
	// 	if previousFileIdentifier.Size == currentFileIdentifier.Size {
	// 		var err error
	// 		if previousFileIdentifier.Hash == "" {
	// 			previousFileIdentifier.Hash, err = getFileHash(previousFileIdentifier.Path)
	// 			if err != nil {
	// 				return "", err
	// 			}
	// 		}
	// 		currentFileIdentifier.Hash, err = getFileHash(currentFileIdentifier.Path)
	// 		if err != nil {
	// 			return "", err
	// 		}
	// 		if previousFileIdentifier.Hash == currentFileIdentifier.Hash {
	// 			if lastAppendedIndex != i-1 {

	// 			}
	// 		}
	// 	}
	// }

	return "", nil
}
