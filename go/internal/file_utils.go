package internal

import (
	"os"
	"path/filepath"
	"time"
)

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

// TODO: use FileFilterMode
func getFilteredFileDetailsFromDirectoryTree(rootFilePath string, fileFilterMode FileFilterMode) ([]FileDetail, error) {
	var fileDetails []FileDetail
	err := filepath.WalkDir(rootFilePath, func(filePath string, dirEntry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fileInfo, err := dirEntry.Info()
		if err != nil {
			return err
		}
		fileDetails = append(fileDetails, FileDetail{
			Path:             filePath,
			ModificationTime: fileInfo.ModTime(),
			Size:             fileInfo.Size(),
			IsDirectory:      dirEntry.IsDir(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fileDetails, nil
}

func joinOutputBasePathWithRelativeInputPath(inputBasePath, inputFullPath, outputBasePath string) (string, error) {
	relativePath, err := filepath.Rel(inputBasePath, inputFullPath)
	if err != nil {
		return "", err
	}
	return filepath.Join(outputBasePath, relativePath), nil
}
