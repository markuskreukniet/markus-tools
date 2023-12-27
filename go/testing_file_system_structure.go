package main

import (
	"os"
	"path/filepath"
)

var (
	directoryEmpty               = "directory empty"
	directory1                   = "directory 1"
	directory2                   = "directory 2"
	directory2WithDirectoryEmpty = filepath.Join(directory2, directoryEmpty)
	directory2WithDirectory3     = filepath.Join(directory2, "directory 3")
	directory2WithDirectory4     = filepath.Join(directory2, "directory 4")

	txtFile1 = filepath.Join(directory1, "file 1.txt")
	txtFile2 = filepath.Join(directory1, "file 2.txt")
	txtFile3 = filepath.Join(directory2WithDirectory3, "file 3.txt")
	txtFile4 = filepath.Join(directory2WithDirectory3, "file 4.txt")
	txtFile5 = filepath.Join(directory2WithDirectory3, "file 5.txt")
	txtFile6 = filepath.Join(directory2WithDirectory4, "file 6.txt")

	emptyPathEndParts []string
)

func testingCreateTempFileSystemStructureOrGetEmptyString(directoryPathEndParts, filePathEndParts []string) (string, error) {
	if len(directoryPathEndParts) == 0 {
		return "", nil
	}
	return testingCreateTempFileSystemStructure(directoryPathEndParts, filePathEndParts)
}

func testingCreateTempFileSystemStructure(directoryPathEndParts, filePathEndParts []string) (string, error) {
	tempDirectory, err := os.MkdirTemp("", "markus-tools go test")
	if err != nil {
		return "", err
	}
	for _, part := range directoryPathEndParts {
		if err := os.MkdirAll(filepath.Join(tempDirectory, part), 0755); err != nil {
			return "", err
		}
	}
	for _, part := range filePathEndParts {
		if err := os.WriteFile(filepath.Join(tempDirectory, part), []byte{}, 0666); err != nil {
			return "", err
		}
	}
	return tempDirectory, nil
}