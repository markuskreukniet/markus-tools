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

func createFileInfoGroupsByHash(files []FileInfo, onlyDuplicates bool) ([][]FileInfo, error) {
	if len(files) == 0 {
		return nil, nil
	}

	type filesByFileSize struct {
		fileSize int64
		files    []FileInfo
	}

	var result [][]FileInfo
	var groups []filesByFileSize
	sizeIndex := 0

	appendGroup := func(file FileInfo) {
		groups = append(groups, filesByFileSize{
			fileSize: file.GetSize(),
			files:    []FileInfo{file},
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].GetSize() < files[j].GetSize()
	})

	appendGroup(files[0])

	for i := 1; i < len(files); i++ {
		if files[i].GetSize() == groups[sizeIndex].files[0].GetSize() {
			groups[sizeIndex].files = append(groups[sizeIndex].files, files[i])
		} else {
			appendGroup(files[i])
			sizeIndex++
		}
	}

	for _, group := range groups {
		if len(group.files) > 1 {
			hashMap := make(map[string][]FileInfo)
			for _, file := range group.files {
				hash, err := CreateFileHash(file.GetAbsolutePath())
				if err != nil {
					return nil, err
				}
				hashMap[hash] = append(hashMap[hash], file)
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

// TODO: does work, but can be improved?
func CreateFileSystemFileByHashGroups(files []FileSystemFile, onlyDuplicates bool) ([][]FileSystemFile, error) {
	if len(files) == 0 {
		return nil, nil
	}

	type filesByFileSize struct {
		fileSize int64
		files    []FileSystemFile
	}

	var result [][]FileSystemFile
	var groups []filesByFileSize
	sizeIndex := 0

	sort.Slice(files, func(i, j int) bool {
		return files[i].FileMetadata.Size < files[j].FileMetadata.Size
	})

	groups = append(groups, filesByFileSize{
		fileSize: files[0].FileMetadata.Size,
		files:    []FileSystemFile{files[0]},
	})

	for i := 1; i < len(files); i++ {
		if files[i].FileMetadata.Size == groups[sizeIndex].files[0].FileMetadata.Size {
			groups[sizeIndex].files = append(groups[sizeIndex].files, files[i])
		} else {
			groups = append(groups, filesByFileSize{
				fileSize: files[i].FileMetadata.Size,
				files:    []FileSystemFile{files[i]},
			})
			sizeIndex++
		}
	}

	for _, group := range groups {
		if len(group.files) > 1 {
			hashMap := make(map[string][]FileSystemFile)
			for _, file := range group.files {
				var err error
				if file.FileMetadata.Hash, err = CreateFileHash(file.FileMetadata.Path); err != nil {
					return nil, err
				}
				hashMap[file.FileMetadata.Hash] = append(hashMap[file.FileMetadata.Hash], file)
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

func CreateFileHash(filePath string) (string, error) {
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
