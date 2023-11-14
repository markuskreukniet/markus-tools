package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"sort"
	"strings"
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

func appendFileIdentifier(fileIdentifiers *[]FileIdentifier, fileDetail FileDetail) {
	*fileIdentifiers = append(*fileIdentifiers, FileIdentifier{
		Path: fileDetail.Path,
		Size: fileDetail.Size,
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

func writeNewlineString(builder *strings.Builder) (int, error) {
	bytesWritten, err := builder.WriteString("\n")
	if err != nil {
		return bytesWritten, err
	}
	return bytesWritten, nil
}

// TODO: use writeNewlineString and cleanup
func duplicateFilesToNewlineSeparatedString(duplicateFiles []DuplicateFile) (string, error) {
	if len(duplicateFiles) == 0 {
		return "", nil
	}
	newLine := "\n"
	var result strings.Builder
	_, err := result.WriteString(duplicateFiles[0].Path)
	if err != nil {
		return "", err
	}
	for i := 1; i < len(duplicateFiles); i++ {
		_, err = result.WriteString(newLine)
		if err != nil {
			return "", err
		}
		if duplicateFiles[i].Hash != duplicateFiles[i-1].Hash {
			_, err = result.WriteString(newLine)
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

func getDuplicateFilesAsNewlineSeparatedString(uniqueFileSystemNodes []FileSystemNode) (string, error) {
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
	var duplicateFiles []DuplicateFile
	var lastAppendedIndex = -1
	for i := 1; i < len(fileIdentifiers); i++ {
		previousIndex := i - 1
		previousFileIdentifier := fileIdentifiers[previousIndex] // TODO: not needed
		currentFileIdentifier := fileIdentifiers[i]              // TODO: not needed
		if previousFileIdentifier.Size == currentFileIdentifier.Size {
			var err error
			if previousFileIdentifier.Hash == "" {
				previousFileIdentifier.Hash, err = getFileHash(previousFileIdentifier.Path)
				if err != nil {
					return "", err
				}
			}
			currentFileIdentifier.Hash, err = getFileHash(currentFileIdentifier.Path)
			if err != nil {
				return "", err
			}
			if previousFileIdentifier.Hash == currentFileIdentifier.Hash {
				if lastAppendedIndex != previousIndex {
					appendDuplicateFile(&duplicateFiles, previousFileIdentifier)
				}
				appendDuplicateFile(&duplicateFiles, currentFileIdentifier)
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
