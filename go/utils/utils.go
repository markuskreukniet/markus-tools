package utils

import (
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

// DuplicateFileGroups
type DuplicateFileGroup struct {
	Identifier string
	FilePaths  []string
}

type DuplicateFileGroups []DuplicateFileGroup

func (groups DuplicateFileGroups) AppendByIdentifier(identifier, filePath string) bool {
	for i, group := range groups {
		if identifier == group.Identifier {
			groups[i].FilePaths = append(groups[i].FilePaths, filePath)
			return true
		}
	}
	return false
}
