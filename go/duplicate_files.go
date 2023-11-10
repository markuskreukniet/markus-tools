package main

type FileSystemNode struct {
	Path        string
	IsDirectory bool
}

type FileIdentifier struct {
	Path string
	Size int64
	Hash string
}

// func duplicateFilesString(uniqueFileSystemNodes []FileSystemNode) (string, error) {
// 	var fileIdentifiers []FileIdentifier
// 	for _, value := range uniqueFileSystemNodes {
// 		if !value.IsDirectory {
// 			fileDetail, err := getFileDetail(value.Path)
// 			if err != nil {
// 				return "", err
// 			}
// 			fileIdentifiers = append(fileIdentifiers, FileIdentifier{
// 				Path: fileDetail.Path,
// 				Size: fileDetail.Size,
// 				Hash: "",
// 			})
// 		} else {

// 		}
// 	}
// }
