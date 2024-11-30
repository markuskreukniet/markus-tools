package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

func plainTextFilesToTextToJSON(uniqueFileSystemNodes []utils.FileSystemNode) string {
	text, err := plainTextFilesToText(uniqueFileSystemNodes)
	if err != nil {
		return errorToJSONFunctionResult(err)
	}
	return resultToJSONFunctionResult(text)
}

func plainTextFilesToText(uniqueFileSystemNodes []utils.FileSystemNode) (string, error) {
	var files []utils.FTextFilesFileInfo

	handler := func(file utils.CompleteFileInfo) {
		files = append(files, utils.FTextFilesFileInfo{
			Name:         file.Name,
			AbsolutePath: file.AbsolutePath,
		})
	}

	for _, node := range uniqueFileSystemNodes {
		if err :=
			utils.WalkFilterAndHandleFileInfo(node, utils.NonZeroByteFiles, utils.TextFiles, handler); err != nil {
			return "", err
		}
	}

	if len(files) == 0 {
		return "", nil
	}

	var result strings.Builder

	addNameAndLines := func(fileInfo utils.FTextFilesFileInfo) error {
		if _, err := result.WriteString(fileInfo.Name); err != nil {
			return err
		}

		file, err := os.Open(fileInfo.AbsolutePath)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			if _, err := utils.WriteNewlineString(&result); err != nil {
				return err
			}
			if _, err = result.WriteString(scanner.Text()); err != nil {
				return err
			}
		}

		return scanner.Err()
	}

	if err := addNameAndLines(files[0]); err != nil {
		return "", err
	}

	for _, file := range files[1:] {
		if _, err := utils.WriteTwoNewlineStrings(&result); err != nil {
			return "", err
		}
		if err := addNameAndLines(file); err != nil {
			return "", err
		}
	}

	return result.String(), nil
}
