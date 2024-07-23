package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"sort"
	"strings"
)

// WriteNewline
func WriteNewlineString(builder *strings.Builder) (int, error) {
	bytesWritten, err := builder.WriteString("\n")
	if err != nil {
		return bytesWritten, err
	}
	return bytesWritten, nil
}

func WriteTwoNewlineStrings(builder *strings.Builder) (int, error) {
	bytesWritten, err := builder.WriteString("\n\n")
	if err != nil {
		return bytesWritten, err
	}
	return bytesWritten, nil
}

// TODO: is this an improved version? Naming in this function is not that good
// func CreateFileHashGroups(files []FileData, onlyDuplicates bool) (FilesDataGroups, error) {
// 	if len(files) == 0 {
// 		return nil, nil
// 	}

// 	// Sort files by their size
// 	sort.Slice(files, func(i, j int) bool {
// 		return files[i].FileMetadata.Size < files[j].FileMetadata.Size
// 	})

// 	var result FilesDataGroups
// 	groupMap := make(map[int64][]FileData)

// 	// Group files by their size
// 	for _, file := range files {
// 		groupMap[file.FileMetadata.Size] = append(groupMap[file.FileMetadata.Size], file)
// 	}

// 	for _, group := range groupMap {
// 		if len(group) > 1 {
// 			hashMap := make(map[string][]FileData)
// 			for _, file := range group {
// 				hash, err := HashFile(file.FileMetadata.Path)
// 				if err != nil {
// 					return nil, err
// 				}
// 				hashMap[hash] = append(hashMap[hash], file)
// 			}
// 			for _, hashedFiles := range hashMap {
// 				if len(hashedFiles) > 1 || !onlyDuplicates {
// 					result = append(result, CreateFilesDataGroup("", hashedFiles))
// 				}
// 			}
// 		} else if !onlyDuplicates {
// 			result = append(result, CreateFilesDataGroup("", group))
// 		}
// 	}

// 	return result, nil
// }

// // TODO: does work, but can be improved with code from above here?
// func CreateFileSystemFileExtraGroups(files []FileSystemFileExtra, onlyDuplicates bool) ([][]FileSystemFileExtra, error) {
// 	if len(files) == 0 {
// 		return nil, nil
// 	}

// 	type filesByFileSize struct {
// 		fileSize int64
// 		files    []FileSystemFileExtra
// 	}

// 	var result [][]FileSystemFileExtra
// 	var groups []filesByFileSize
// 	sizeIndex := 0

// 	sort.Slice(files, func(i, j int) bool {
// 		return files[i].FileSystemFile.FileMetadata.Size < files[j].FileSystemFile.FileMetadata.Size
// 	})

// 	groups = append(groups, filesByFileSize{
// 		fileSize: files[0].FileSystemFile.FileMetadata.Size,
// 		files:    []FileSystemFileExtra{files[0]},
// 	})

// 	for i := 1; i < len(files); i++ {
// 		if files[i].FileSystemFile.FileMetadata.Size == groups[sizeIndex].files[0].FileSystemFile.FileMetadata.Size {
// 			groups[sizeIndex].files = append(groups[sizeIndex].files, files[i])
// 		} else {
// 			groups = append(groups, filesByFileSize{
// 				fileSize: files[i].FileSystemFile.FileMetadata.Size,
// 				files:    []FileSystemFileExtra{files[i]},
// 			})
// 			sizeIndex++
// 		}
// 	}

// 	for _, group := range groups {
// 		if len(group.files) > 1 {
// 			hashMap := make(map[string][]FileSystemFileExtra)
// 			for _, file := range group.files {
// 				var err error
// 				if file.Hash, err = HashFile(file.FileSystemFile.Path); err != nil {
// 					return nil, err
// 				}
// 				hashMap[file.Hash] = append(hashMap[file.Hash], file)
// 			}
// 			for _, hashedFiles := range hashMap {
// 				if len(hashedFiles) > 1 || !onlyDuplicates {
// 					result = append(result, hashedFiles)
// 				}
// 			}
// 		} else if !onlyDuplicates {
// 			result = append(result, group.files)
// 		}
// 	}

// 	return result, nil
// }

// TODO: does work, but can be improved with code from above here?
func CreateFileHashGroups(files []FileData, onlyDuplicates bool) (FilesDataGroups, error) {
	if len(files) == 0 {
		return nil, nil
	}

	type filesByFileSize struct {
		fileSize int64
		files    []FileData
	}

	var result FilesDataGroups
	var groups []filesByFileSize
	sizeIndex := 0

	sort.Slice(files, func(i, j int) bool {
		return files[i].FileMetadata.Size < files[j].FileMetadata.Size
	})

	groups = append(groups, filesByFileSize{
		fileSize: files[0].FileMetadata.Size,
		files:    []FileData{files[0]},
	})

	for i := 1; i < len(files); i++ {
		if files[i].FileMetadata.Size == groups[sizeIndex].files[0].FileMetadata.Size {
			groups[sizeIndex].files = append(groups[sizeIndex].files, files[i])
		} else {
			groups = append(groups, filesByFileSize{
				fileSize: files[i].FileMetadata.Size,
				files:    []FileData{files[i]},
			})
			sizeIndex++
		}
	}

	for _, group := range groups {
		if len(group.files) > 1 {
			hashMap := make(map[string][]FileData)
			for _, file := range group.files {
				var err error
				if file.Identifier, err = HashFile(file.FileMetadata.Path); err != nil {
					return nil, err
				}
				hashMap[file.Identifier] = append(hashMap[file.Identifier], file)
			}
			for _, hashedFiles := range hashMap {
				if len(hashedFiles) > 1 || !onlyDuplicates {
					result = append(result, CreateFilesDataGroup("", hashedFiles))
				}
			}
		} else if !onlyDuplicates {
			result = append(result, CreateFilesDataGroup("", group.files))
		}
	}

	return result, nil
}

func CreateDuplicateFileGroups(files []FileData) (FilesDataGroups, error) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].FileMetadata.Size < files[j].FileMetadata.Size
	})
	var groups FilesDataGroups
	for i := 1; i < len(files); i++ {
		previousIndex := i - 1
		if files[previousIndex].FileMetadata.Size == files[i].FileMetadata.Size {
			var err error
			if files[previousIndex].Identifier == "" {
				if files[previousIndex].Identifier, err = HashFile(files[previousIndex].FileMetadata.Path); err != nil {
					return nil, err
				}
			}
			if files[i].Identifier, err = HashFile(files[i].FileMetadata.Path); err != nil {
				return nil, err
			}
			if !groups.DidAppendByFileDataIdentifier(files[i]) {
				for j := 0; j <= previousIndex; j++ {
					if files[i].Identifier == files[j].Identifier {
						groups = append(groups, CreateFilesDataGroup(files[i].Identifier, []FileData{files[j], files[i]}))
						break
					}
				}
			}
		}
	}
	return groups, nil
}

func HashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// sha256 is generally faster and more secure than SHA1
	// SHA1 is generally faster and more secure than MD5
	hashGenerator := sha256.New()
	if _, err := io.Copy(hashGenerator, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hashGenerator.Sum(nil)), nil
}
