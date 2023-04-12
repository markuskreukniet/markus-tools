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

func WriteTwoNewlineStrings(writer io.Writer) (int, error) {
	bytesWritten, err := io.WriteString(writer, "\n\n")
	if err != nil {
		return bytesWritten, err
	}
	return bytesWritten, nil
}

func CreateDuplicateFileInfoGroupsByHash[T DuplicateFileInfo](files []T, onlyDuplicates bool) ([][]T, error) {
	if len(files) == 0 {
		return nil, nil
	}

	type filesByFileSize struct {
		fileSize int64
		files    []T
	}

	var result [][]T
	var groups []filesByFileSize
	groupIndex := 0

	appendGroup := func(file T) {
		groups = append(groups, filesByFileSize{
			fileSize: file.GetSize(),
			files:    []T{file},
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].GetSize() < files[j].GetSize()
	})

	appendGroup(files[0])

	for _, file := range files[1:] {
		if file.GetSize() == groups[groupIndex].files[0].GetSize() {
			groups[groupIndex].files = append(groups[groupIndex].files, file)
		} else {
			appendGroup(file)
			groupIndex++
		}
	}

	for _, group := range groups {
		if len(group.files) > 1 {
			hashMap := make(map[string][]T)
			for _, file := range group.files {
				hash, err := CreateFileHash(file.GetPath())
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

func IsBlank(s string) bool {
	return s == ""
}
