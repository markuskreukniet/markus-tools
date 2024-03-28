package utils

import (
	"io"
	"os"
	"path/filepath"
	"time"
	"unicode"
)

// TODO: should be timeModified instead of ModificationTime. Also search in project on ModificationTime for more renaming
type FileDetail struct {
	Path             string
	ModificationTime time.Time
	Size             int64
}

type FileSystemNode struct {
	Path        string
	IsDirectory bool
}

type (
	fileFilterMode int
	fileType       int
)

const (
	files fileFilterMode = iota
	FilesWithoutZeroByteFiles
	FilesAndDirectories
	filesAndDirectoriesWithoutZeroByteFiles
	Directories
)

const (
	AllFiles fileType = iota
	PlainTextFiles
)

func IsNonZeroByteFileATextFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read the first 512 or less to check for non-text characters. DetectContentType of package 'net/http' works with a similar check.
	bytes := make([]byte, 512)
	numberOfBytesRead, err := file.Read(bytes)
	if err != nil && err != io.EOF {
		return false, err
	}
	for _, byte := range bytes[:numberOfBytesRead] {
		if !unicode.IsPrint(rune(byte)) && !unicode.IsSpace(rune(byte)) {
			return false, nil
		}
	}
	return true, nil
}

func IsFileDetailNonZeroByte(detail FileDetail) bool {
	if detail.Size > 0 {
		return true
	}
	return false
}

func GetFileDetail(filePath string) (FileDetail, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return FileDetail{}, err
	}
	return FileDetail{
		Path:             filePath,
		ModificationTime: fileInfo.ModTime(),
		Size:             fileInfo.Size(),
	}, nil
}

func WalkFileDetails(rootFilePath string, mode fileFilterMode, fileType fileType, handler func(FileDetail)) error {
	return filepath.Walk(rootFilePath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		size := fileInfo.Size()
		isDir := fileInfo.IsDir()

		// is file check
		if !isDir && mode == Directories {
			return nil
		}

		// is directory check
		if isDir && (mode == files || mode == FilesWithoutZeroByteFiles) {
			return nil
		}

		// zero byte check
		if size == 0 && (mode == FilesWithoutZeroByteFiles || mode == filesAndDirectoriesWithoutZeroByteFiles) {
			return nil
		}

		// file type check
		if fileType == PlainTextFiles {
			isTextFile, err := IsNonZeroByteFileATextFile(filePath)
			if err != nil || !isTextFile {
				return err
			}
		}

		handler(FileDetail{
			Path:             filePath,
			ModificationTime: fileInfo.ModTime(),
			Size:             size,
		})
		return nil
	})
}

func AppendFileDetails(appendFileDetail func(detail FileDetail), uniqueFileSystemNodes []FileSystemNode, mode fileFilterMode) error {
	for _, node := range uniqueFileSystemNodes {
		if node.IsDirectory {
			if err := WalkFileDetails(node.Path, mode, AllFiles, func(detail FileDetail) {
				appendFileDetail(detail)
			}); err != nil {
				return err
			}
		} else {
			detail, err := GetFileDetail(node.Path)
			if err != nil {
				return err
			}
			if IsFileDetailNonZeroByte(detail) {
				appendFileDetail(detail)
			}
		}
	}
	return nil
}

func FileOrDirectoryExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
