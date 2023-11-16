package main

import (
	"os"
	"path/filepath"
)

func synchronizeDirectoryTreesToJSON(sourceDirectory, destinationDirectory string) string {
	functionResult := FunctionResult{
		Result:       "",
		ErrorMessage: "",
	}
	err := synchronizeDirectoryTrees(sourceDirectory, destinationDirectory)
	if err != nil {
		functionResult.ErrorMessage = err.Error()
	}
	return jsonMarshalWithFallbackJSONError(functionResult)
}

func synchronizeDirectoryTrees(sourceDirectory, destinationDirectory string) error {
	destinationFilePathModificationTimeMap, err := getFilePathModificationTimeMapFromDirectoryTree(destinationDirectory)
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
	})
	if err != nil {
		return err
	}
	for key := range destinationFilePathModificationTimeMap {
		err := os.RemoveAll(key)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return err
}
