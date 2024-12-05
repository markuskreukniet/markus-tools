package main

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

func synchronizeDirectoryTreesToJSON(sourceDirectory, destinationDirectory string) string {
	if err := synchronizeDirectoryTrees(sourceDirectory, destinationDirectory); err != nil {
		return errorToJSONFunctionResult(err)
	}
	return defaultJSONFunctionResult()
}

func joinOutputBasePathWithRelativeInputPath(inputBasePath, inputFullPath, outputBasePath string) (string, error) {
	relativePath, err := filepath.Rel(inputBasePath, inputFullPath)
	if err != nil {
		return "", err
	}
	return filepath.Join(outputBasePath, relativePath), nil
}

func getFilePathModificationTimeMapFromDirectoryTree(directoryPath string) (map[string]time.Time, error) {
	filePathModificationTimeMap := make(map[string]time.Time)

	handler := func(file utils.CompleteFileInfo) {
		filePathModificationTimeMap[file.AbsolutePath] = file.TimeModified
	}

	err := utils.WalkFilterAndHandleFileInfoDirectory(directoryPath, utils.FilesAndDirectories, utils.AllFiles, handler)
	return filePathModificationTimeMap, err
}

// Copying files in this function could be faster with buffering.
// However, to determine an optimal buffer size for copying a file,
// we need to know the block size of the storage device.
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
	if _, err = io.Copy(destinationFile, sourceFile); err != nil {
		return err
	}
	return os.Chmod(destinationFilePath, fileMode)
}

func synchronizeDirectoryTrees(sourceDirectory, destinationDirectory string) error {
	destinationFilePathModificationTimeMap, err := getFilePathModificationTimeMapFromDirectoryTree(destinationDirectory)
	if err != nil {
		return err
	}

	if err = filepath.Walk(sourceDirectory, func(sourceFilePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		destinationFilePath, err := joinOutputBasePathWithRelativeInputPath(sourceDirectory, sourceFilePath, destinationDirectory)
		if err != nil {
			return err
		}
		isDir := fileInfo.IsDir()
		value, ok := destinationFilePathModificationTimeMap[destinationFilePath]
		if !isDir && (!ok || fileInfo.ModTime().After(value)) {
			err = copyFileWithFileMode(sourceFilePath, destinationFilePath, fileInfo.Mode())
		} else if isDir && !ok {
			err = os.Mkdir(destinationFilePath, fileInfo.Mode())
		}
		if ok {
			delete(destinationFilePathModificationTimeMap, destinationFilePath)
		}
		return err
	}); err != nil {
		return err
	}
	for key := range destinationFilePathModificationTimeMap {
		if err := os.RemoveAll(key); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return err
}
