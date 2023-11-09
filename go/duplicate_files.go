package main

import "time"

type FileSystemNode struct {
	Path        string
	IsDirectory bool
}

type FileMetadata struct {
	ModificationTime time.Time
	Size             int64
	Hash             string
}

// func duplicateFiles(uniqueFileSystemNodes []FileSystemNode) {
// 	var filesMetadata []FileMetadata
// 	for _, value := range uniqueFileSystemNodes {
// 		if !value.IsDirectory {
// 			//
// 		}
// 	}
// }
