package main

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// TODO: FileDetail has maybe unused fields
type FileDetail struct {
	Path             string
	ModificationTime time.Time
	Size             int64
	IsDirectory      bool
}

type FileFilterMode int

const (
	files FileFilterMode = iota
	filesWithoutZeroByteFiles
	filesAndDirectories
	filesAndDirectoriesWithoutZeroByteFiles
	directories
)

func getFileDetail(filePath string) (FileDetail, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return FileDetail{}, err
	}
	return FileDetail{
		Path:             filePath,
		ModificationTime: fileInfo.ModTime(),
		Size:             fileInfo.Size(),
		IsDirectory:      fileInfo.IsDir(),
	}, nil
}

// WIP
// func getAsdf(filePath string, directoryTree bool, fileFilterMode FileFilterMode) ([]FileDetail, error) {
// 	var fileDetails []FileDetail
// 	var stack []string
// 	for stackLength := len(stack); stackLength > 0; stackLength = len(stack) {
// 		var stackLengthMinOne int = stackLength - 1
// 		fileDetail := stack[stackLengthMinOne]
// 		stack = stack[:stackLengthMinOne]

// 		// read files from directory
// 	}
// }

// TODO: func useless?
func escapedStringToJSON(escapedString string) string {
	return `"` + escapedString + `"`
}

func jsonMarshalWithFallbackJSONError(nonJSON string) string {
	jsonBytes, err := json.Marshal(nonJSON)
	if err != nil {
		return escapedStringToJSON("json.Marshal error")
	}
	return string(jsonBytes)
}

// Copying files in this function could be faster with buffering.
// However, to determine an optimal buffer size for copying a file, we need to know the block size of the storage device.
// Determining such block sizes is relatively hard with only official Go packages.
func copyFileWithFileMode(sourceFilePath string, destinationFilePath string, fileMode fs.FileMode) error {
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	destinationFile, err := os.Create(destinationFilePath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}
	return os.Chmod(destinationFilePath, fileMode)
}

func getFilePathModificationTimeMapFromDirectoryTree(rootFilePath string) (map[string]time.Time, error) {
	filePathModificationTimeMap := make(map[string]time.Time)
	err := walkFileDetails(rootFilePath, filesAndDirectories, func(fileDetail FileDetail) {
		filePathModificationTimeMap[fileDetail.Path] = fileDetail.ModificationTime
	})
	return filePathModificationTimeMap, err
}

// TODO: not used, but does get tested
func getFilteredFileDetailsSliceFromDirectoryTree(rootFilePath string, fileFilterMode FileFilterMode) ([]FileDetail, error) {
	var fileDetails []FileDetail
	err := walkFileDetails(rootFilePath, fileFilterMode, func(fileDetail FileDetail) {
		fileDetails = append(fileDetails, fileDetail)
	})
	return fileDetails, err
}

func walkFileDetails(rootFilePath string, fileFilterMode FileFilterMode, handler func(FileDetail)) error {
	return filepath.Walk(rootFilePath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		size := fileInfo.Size()
		isDir := fileInfo.IsDir()

		// is file check
		if !isDir && fileFilterMode == directories {
			return nil
		}

		// is directory check
		if isDir && (fileFilterMode == files || fileFilterMode == filesWithoutZeroByteFiles) {
			return nil
		}

		// zero byte check
		if size == 0 && (fileFilterMode == filesWithoutZeroByteFiles || fileFilterMode == filesAndDirectoriesWithoutZeroByteFiles) {
			return nil
		}

		handler(FileDetail{
			Path:             filePath,
			ModificationTime: fileInfo.ModTime(),
			Size:             size,
			IsDirectory:      isDir,
		})
		return nil
	})
}

func joinOutputBasePathWithRelativeInputPath(inputBasePath, inputFullPath, outputBasePath string) (string, error) {
	relativePath, err := filepath.Rel(inputBasePath, inputFullPath)
	if err != nil {
		return "", err
	}
	return filepath.Join(outputBasePath, relativePath), nil
}