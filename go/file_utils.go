package main

import (
	"os"
	"path/filepath"
	"time"
)

// TODO: some things should be in other files since they are only used in one file
// TODO: FileDetail has maybe unused fields
type FileDetail struct {
	Path             string
	ModificationTime time.Time
	Size             int64
	IsDirectory      bool
}

type (
	FileFilterMode int
	FileType       int
)

const (
	files FileFilterMode = iota
	filesWithoutZeroByteFiles
	filesAndDirectories
	filesAndDirectoriesWithoutZeroByteFiles
	directories
)

const (
	allFiles FileType = iota
	plainTextFiles
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

// TODO: not used, but does get tested
func getFilteredFileDetailsSliceFromDirectoryTree(rootFilePath string, fileFilterMode FileFilterMode, fileType FileType) ([]FileDetail, error) {
	var fileDetails []FileDetail
	err := walkFileDetails(rootFilePath, fileFilterMode, fileType, func(fileDetail FileDetail) {
		fileDetails = append(fileDetails, fileDetail)
	})
	return fileDetails, err
}

func walkFileDetails(rootFilePath string, fileFilterMode FileFilterMode, fileType FileType, handler func(FileDetail)) error {
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

		// file type check
		if fileType == plainTextFiles {
			isTextFile, err := isZeroByteFileATextFile(filePath)
			if err != nil || !isTextFile {
				return err
			}
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
