package main

import "sort"

type FileSystemNode struct {
	Path        string
	IsDirectory bool
}

type FileIdentifier struct {
	Path string
	Size int64
	Hash string
}

func appendFileIdentifier(fileIdentifiers *[]FileIdentifier, fileDetail FileDetail) []FileIdentifier {
	return append(*fileIdentifiers, FileIdentifier{
		Path: fileDetail.Path,
		Size: fileDetail.Size,
		Hash: "",
	})
}

func duplicateFilesString(uniqueFileSystemNodes []FileSystemNode) (string, error) {
	var fileIdentifiers []FileIdentifier
	for _, value := range uniqueFileSystemNodes {
		if !value.IsDirectory {
			fileDetail, err := getFileDetail(value.Path)
			if err != nil {
				return "", err
			}
			if fileDetail.Size > 0 {
				appendFileIdentifier(&fileIdentifiers, fileDetail)
			}
		} else {
			err := walkFileDetails(value.Path, filesAndDirectoriesWithoutZeroByteFiles, func(fileDetail FileDetail) {
				appendFileIdentifier(&fileIdentifiers, fileDetail)
			})
			if err != nil {
				return "", err
			}
		}
	}
	sort.Slice(fileIdentifiers, func(i, j int) bool {
		return fileIdentifiers[i].Size < fileIdentifiers[j].Size
	})
	var duplicates []FileIdentifier
	var lastAppendedIndex = -1
	for i := 1; i < len(fileIdentifiers); i++ {
		previousFileIdentifier := fileIdentifiers[i-1]
		currentFileIdentifier := fileIdentifiers[i]
		if previousFileIdentifier.Size == currentFileIdentifier.Size {

		}
	}
	//
	return "", nil
}
