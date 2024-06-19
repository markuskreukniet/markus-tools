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
				if files[previousIndex].Identifier, err = getFileHash(files[previousIndex].FileMetadata.FilePath); err != nil {
					return nil, err
				}
			}
			if files[i].Identifier, err = getFileHash(files[i].FileMetadata.FilePath); err != nil {
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

func getFileHash(filePath string) (string, error) {
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
