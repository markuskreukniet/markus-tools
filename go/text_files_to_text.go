package main

import (
	"os"
	"unicode"
)

func isZeroByteFileATextFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read the first 512 bytes to check for non-text characters. DetectContentType of package 'net/http' works with a similar check.
	bytes := make([]byte, 512)
	_, err = file.Read(bytes)
	if err != nil {
		return false, err
	}
	for _, byte := range bytes {
		if !unicode.IsPrint(rune(byte)) && !unicode.IsSpace(rune(byte)) {
			return false, nil
		}
	}
	return true, nil
}

// func textFilesToText(uniqueFileSystemNodes []FileSystemNode) (string, error) {

// }
