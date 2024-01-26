package main

import (
	"os"
	"path/filepath"
	"time"
)

type fileDetail struct {
	Path             string
	ModificationTime time.Time
	Size             int64
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

func isFileDetailNonZeroByte(detail fileDetail) bool {
	if detail.Size > 0 {
		return true
	}
	return false
}

func getFileDetail(filePath string) (fileDetail, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fileDetail{}, err
	}
	return fileDetail{
		Path:             filePath,
		ModificationTime: fileInfo.ModTime(),
		Size:             fileInfo.Size(),
	}, nil
}

func walkFileDetails(rootFilePath string, fileFilterMode FileFilterMode, fileType FileType, handler func(fileDetail)) error {
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
			isTextFile, err := isNonZeroByteFileATextFile(filePath)
			if err != nil || !isTextFile {
				return err
			}
		}

		handler(fileDetail{
			Path:             filePath,
			ModificationTime: fileInfo.ModTime(),
			Size:             size,
		})
		return nil
	})
}
