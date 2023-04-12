package utils

import (
	"io"
	"os"
	"path/filepath"
	"time"
	"unicode"
)

// The F prefix of a struct name means feature.

type DuplicateFileInfo interface {
	GetSize() int64
	GetPath() string
}

// FDateRangeFileInfo implements DuplicateFileInfo
type FDateRangeFileInfo struct {
	Size         int64
	Path         string
	Name         string
	TimeModified time.Time
}

func (info FDateRangeFileInfo) GetSize() int64 {
	return info.Size
}

func (info FDateRangeFileInfo) GetPath() string {
	return info.Path
}

// FDuplicateFilesFileInfo implements DuplicateFileInfo
type FDuplicateFilesFileInfo struct {
	Size int64
	Path string
}

func (info FDuplicateFilesFileInfo) GetSize() int64 {
	return info.Size
}

func (info FDuplicateFilesFileInfo) GetPath() string {
	return info.Path
}

type FTextFilesFileInfo struct {
	Name         string
	AbsolutePath string
}

type CompleteFileInfo struct {
	Name                  string
	AbsoluteDirectoryPath string
	AbsolutePath          string
	TimeModified          time.Time
	Size                  int64
	IsDirectory           bool
}

type FileData struct {
	Content          string
	CompleteFileInfo CompleteFileInfo
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
	NonZeroByteFiles
	FilesAndDirectories
	NonZeroByteFilesAndDirectories
	Directories
)

const (
	AllFiles fileType = iota
	TextFiles
)

const FilePathSeparator = string(filepath.Separator)

func CreateDirectory(filePath string) error {
	return os.Mkdir(filePath, 0755)
}

func IsTextFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read the first 512 or less to check for non-text characters.
	// DetectContentType of package 'net/http' works with a similar check.
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

func resolveDirectoryPath(filePath string, isDirectory bool) string {
	directoryPath := filePath

	if !isDirectory {
		directoryPath = filepath.Dir(filePath)
	}

	return directoryPath
}

func FilterAndHandleFileInfo(
	info os.FileInfo, mode fileFilterMode, fileType fileType, absoluteFilePath string, handler func(CompleteFileInfo),
) error {
	isDir, isRegularFile := info.IsDir(), info.Mode().IsRegular()

	var size int64
	if isRegularFile {
		size = info.Size()
	}

	// is file check
	if isRegularFile && mode == Directories {
		return nil
	}

	// is directory check
	if isDir && (mode == files || mode == NonZeroByteFiles) {
		return nil
	}

	// is zero byte file check
	if isRegularFile && size == 0 && (mode == NonZeroByteFiles || mode == NonZeroByteFilesAndDirectories) {
		return nil
	}

	// is text file check
	if fileType == TextFiles {
		isTextFile, err := IsTextFile(absoluteFilePath)
		if err != nil || !isTextFile {
			return err
		}
	}

	handler(CompleteFileInfo{
		Name:                  info.Name(),
		AbsoluteDirectoryPath: resolveDirectoryPath(absoluteFilePath, isDir),
		AbsolutePath:          absoluteFilePath,
		TimeModified:          info.ModTime(),
		Size:                  size,
		IsDirectory:           isDir,
	})

	return nil
}

func WalkFilterAndHandleFileInfoDirectory(
	filePath string, mode fileFilterMode, fileType fileType, handler func(CompleteFileInfo)) error {
	return filepath.Walk(filePath, func(absoluteFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return FilterAndHandleFileInfo(info, mode, fileType, absoluteFilePath, handler)
	})
}

func WalkFilterAndHandleFileInfo(
	node FileSystemNode, mode fileFilterMode, fileType fileType, handler func(CompleteFileInfo),
) error {
	if node.IsDirectory {
		return WalkFilterAndHandleFileInfoDirectory(node.Path, mode, fileType, handler)
	} else {
		info, err := os.Stat(node.Path)
		if err != nil {
			return err
		}
		return FilterAndHandleFileInfo(info, mode, fileType, node.Path, handler)
	}
}

func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
