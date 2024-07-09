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
			if !groups.AppendByFileDataIdentifier(files[i]) {
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

// TODO: WIP
// func CreateHashFileGroups(files []FileData, onlyDuplicates bool) (FilesDataGroups, error) {
// 	sort.Slice(files, func(i, j int) bool {
// 		return files[i].FileMetadata.Size < files[j].FileMetadata.Size
// 	})

// 	var groups FilesDataGroups
// 	length := len(files)

// 	for i := 0; i < length-1; i++ {
// 		nextIndex := i - 1
// 		if files[i].FileMetadata.Size == files[nextIndex].FileMetadata.Size {
// 			var err error
// 			if files[i].Identifier == "" {
// 				if files[i].Identifier, err = HashFile(files[i].FileMetadata.Path); err != nil {
// 					return nil, err
// 				}
// 			}
// 			if files[nextIndex].Identifier, err = HashFile(files[nextIndex].FileMetadata.Path); err != nil {
// 				return nil, err
// 			}
// 			if !groups.AppendByFileDataIdentifier(files[nextIndex]) {
// 				for j := 0; j <= nextIndex; j++ {
// 					if files[j].Identifier == files[nextIndex].Identifier {
// 						groups = append(groups, CreateFilesDataGroup(files[nextIndex].Identifier, []FileData{files[j], files[nextIndex]}))
// 						break
// 					}
// 				}
// 			}
// 		} else if !onlyDuplicates {
// 			previousIndex := i - 1
// 			if i == length-1 ||
// 				(files[previousIndex].Identifier != files[i].Identifier ||
// 					(files[previousIndex].Identifier == "" && files[i].Identifier == "")) {
// 				groups = append(groups, CreateFilesDataGroup(files[i].Identifier, []FileData{files[i]}))
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
