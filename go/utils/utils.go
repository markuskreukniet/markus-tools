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

// TODO: does work, but can be improved?
func CreateFileSystemFileExtraByHashGroups(files []FileSystemFileExtra, onlyDuplicates bool) ([][]FileSystemFileExtra, error) {
	if len(files) == 0 {
		return nil, nil
	}

	type filesByFileSize struct {
		fileSize int64
		files    []FileSystemFileExtra
	}

	var result [][]FileSystemFileExtra
	var groups []filesByFileSize
	sizeIndex := 0

	sort.Slice(files, func(i, j int) bool {
		return files[i].FileSystemFile.FileMetadata.Size < files[j].FileSystemFile.FileMetadata.Size
	})

	groups = append(groups, filesByFileSize{
		fileSize: files[0].FileSystemFile.FileMetadata.Size,
		files:    []FileSystemFileExtra{files[0]},
	})

	for i := 1; i < len(files); i++ {
		if files[i].FileSystemFile.FileMetadata.Size == groups[sizeIndex].files[0].FileSystemFile.FileMetadata.Size {
			groups[sizeIndex].files = append(groups[sizeIndex].files, files[i])
		} else {
			groups = append(groups, filesByFileSize{
				fileSize: files[i].FileSystemFile.FileMetadata.Size,
				files:    []FileSystemFileExtra{files[i]},
			})
			sizeIndex++
		}
	}

	for _, group := range groups {
		if len(group.files) > 1 {
			hashMap := make(map[string][]FileSystemFileExtra)
			for _, file := range group.files {
				var err error
				if file.Hash, err = HashFile(file.FileSystemFile.FileMetadata.Path); err != nil {
					return nil, err
				}
				hashMap[file.Hash] = append(hashMap[file.Hash], file)
			}
			for _, hashedFiles := range hashMap {
				if len(hashedFiles) > 1 || !onlyDuplicates {
					result = append(result, hashedFiles)
				}
			}
		} else if !onlyDuplicates {
			result = append(result, group.files)
		}
	}

	return result, nil
}

// func CreateDuplicateFileGroups(files []FileData) (FilesDataGroups, error) {
// 	sort.Slice(files, func(i, j int) bool {
// 		return files[i].FileMetadata.Size < files[j].FileMetadata.Size
// 	})
// 	var groups FilesDataGroups
// 	for i := 1; i < len(files); i++ {
// 		previousIndex := i - 1
// 		if files[previousIndex].FileMetadata.Size == files[i].FileMetadata.Size {
// 			var err error
// 			if files[previousIndex].Identifier == "" {
// 				if files[previousIndex].Identifier, err = HashFile(files[previousIndex].FileMetadata.Path); err != nil {
// 					return nil, err
// 				}
// 			}
// 			if files[i].Identifier, err = HashFile(files[i].FileMetadata.Path); err != nil {
// 				return nil, err
// 			}
// 			if !groups.DidAppendByFileDataIdentifier(files[i]) {
// 				for j := 0; j <= previousIndex; j++ {
// 					if files[i].Identifier == files[j].Identifier {
// 						groups = append(groups, CreateFilesDataGroup(files[i].Identifier, []FileData{files[j], files[i]}))
// 						break
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return groups, nil
// }

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
