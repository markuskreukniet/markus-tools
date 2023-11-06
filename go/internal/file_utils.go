package internal

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// TODO: FileDetailMapValue should be part of FileDetail
// TODO: FileDetail and FileDetailMapValue have maybe unused fields
// TODO: add new struct for synchronizeDirectoryTrees to use less memory
type FileDetail struct {
	Path             string
	ModificationTime time.Time
	Size             int64
	IsDirectory      bool
}

type FileDetailMapValue struct {
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

func synchronizeDirectoryTrees(sourceDirectory, destinationDirectory string) error {
	destinationFileDetails, err := getFilteredFileDetailsMapFromDirectoryTree(destinationDirectory, filesAndDirectories)
	if err != nil {
		return err
	}
	err = filepath.Walk(sourceDirectory, func(sourceFilePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		destinationFilePath, err := joinOutputBasePathWithRelativeInputPath(sourceDirectory, sourceFilePath, destinationDirectory)
		if err != nil {
			return err
		}
		isDir := fileInfo.IsDir()
		value, ok := destinationFileDetails[destinationFilePath]
		if !isDir && (!ok || fileInfo.ModTime().After(value.ModificationTime)) {
			err = copyFileWithFileMode(sourceFilePath, destinationFilePath, fileInfo.Mode())
		} else if isDir && !ok {
			err = os.Mkdir(destinationFilePath, fileInfo.Mode())
		}
		if ok {
			delete(destinationFileDetails, destinationFilePath)
		}
		return err
	})
	if err != nil {
		return err
	}
	for key, _ := range destinationFileDetails {
		err := os.RemoveAll(key)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return err
}

// TODO: comment about buffering
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

func getFilteredFileDetailsMapFromDirectoryTree(rootFilePath string, fileFilterMode FileFilterMode) (map[string]FileDetailMapValue, error) {
	fileInfos := make(map[string]FileDetailMapValue)
	err := walkFileDetails(rootFilePath, fileFilterMode, func(fileDetail FileDetail) {
		fileInfos[fileDetail.Path] = FileDetailMapValue{
			ModificationTime: fileDetail.ModificationTime,
			Size:             fileDetail.Size,
			IsDirectory:      fileDetail.IsDirectory,
		}
	})
	return fileInfos, err
}

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
