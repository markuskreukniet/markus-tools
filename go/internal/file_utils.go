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

func getFileDetail(path string) (FileDetail, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return FileDetail{}, err
	}
	return FileDetail{
		Path:             path,
		ModificationTime: fileInfo.ModTime(),
		Size:             fileInfo.Size(),
		IsDirectory:      fileInfo.IsDir(),
	}, nil
}

func joinOutputBasePathWithRelativeInputPath(inputBasePath, inputFullPath, outputBasePath string) (string, error) {
	relativePath, err := filepath.Rel(inputBasePath, inputFullPath)
	if err != nil {
		return "", err
	}
	return filepath.Join(outputBasePath, relativePath), nil
}
