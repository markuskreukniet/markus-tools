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

// DuplicateFilesFileInfo implements DuplicateFileInfo
type DuplicateFilesFileInfo struct {
	Size int64
	Path string
}

func (info DuplicateFilesFileInfo) GetSize() int64 {
	return info.Size
}

func (info DuplicateFilesFileInfo) GetPath() string {
	return info.Path
}

// CompleteFileInfo implements FileInfo
type CompleteFileInfo struct {
	name                  string
	absoluteDirectoryPath string
	absolutePath          string
	timeModified          time.Time
	size                  int64
	isDirectory           bool
}

func (info CompleteFileInfo) GetSize() int64 {
	return info.size
}

func (info CompleteFileInfo) GetPath() string {
	return info.absolutePath
}

type FileSystemFile struct {
	Data         string
	FileMetadata FileMetadata
}

func CreateFileSystemFile(data string, metadata FileMetadata) FileSystemFile {
	return FileSystemFile{
		Data:         data,
		FileMetadata: metadata,
	}
}

type FileMetadata struct {
	Name, DirectoryPath, Path, Hash string
	TimeModified                    time.Time
	Size                            int64
	IsDirectory                     bool // It should be a file type, but there is no use case.
}

func CreateFileMetadata(name, directoryPath, path, hash string, timeModified time.Time, size int64, isDirectory bool) FileMetadata {
	return FileMetadata{
		Name:          name,
		DirectoryPath: directoryPath, // TODO: absoluteDirectoryPath better naming?
		Path:          path,          // TODO: absolutePath better naming?
		TimeModified:  timeModified,
		Size:          size,
		IsDirectory:   isDirectory,
		Hash:          hash,
	}
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

func resolveDirectoryPath(filePath string, isDirectory bool) string {
	directoryPath := filePath

	if !isDirectory {
		directoryPath = filepath.Dir(filePath)
	}

	return directoryPath
}

func ToFileSystemFile(filePath string) (FileSystemFile, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return FileSystemFile{}, err
	}

	isDirectory := info.IsDir()

	return CreateFileSystemFile("",
		CreateFileMetadata(info.Name(), resolveDirectoryPath(filePath, isDirectory), filePath, "", info.ModTime(), info.Size(), isDirectory)), nil
}

func FilterAndHandleFileInfo(
	info os.FileInfo, mode fileFilterMode, fileType fileType, absoluteFilePath string, handler func(DuplicateFileInfo),
) error {
	isDir := info.IsDir()
	isRegularFile := info.Mode().IsRegular()

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
		name:                  info.Name(),
		absoluteDirectoryPath: resolveDirectoryPath(absoluteFilePath, isDir),
		absolutePath:          absoluteFilePath,
		timeModified:          info.ModTime(),
		size:                  size,
		isDirectory:           isDir,
	})

	return nil
}

func WalkFilterAndHandleFileInfo(
	node FileSystemNode, mode fileFilterMode, fileType fileType, handler func(DuplicateFileInfo),
) error {
	if node.IsDirectory {
		return filepath.Walk(node.Path, func(absoluteFilePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			return FilterAndHandleFileInfo(info, mode, fileType, absoluteFilePath, handler)
		})
	} else {
		info, err := os.Stat(node.Path)
		if err != nil {
			return err
		}
		return FilterAndHandleFileInfo(info, mode, fileType, node.Path, handler)
	}
}

func WalkFilterAndHandleFileSystemFile(rootFilePath string, mode fileFilterMode, fileType fileType, handler func(FileSystemFile) error) error {
	return filepath.Walk(rootFilePath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		isDir := fileInfo.IsDir()

		var size int64
		if !isDir {
			size = fileInfo.Size()
		}

		// is file check
		if !isDir && mode == Directories {
			return nil
		}

		// is directory check
		if isDir && (mode == files || mode == NonZeroByteFiles) {
			return nil
		}

		// is zero byte file check
		if !isDir && size == 0 && (mode == NonZeroByteFiles || mode == NonZeroByteFilesAndDirectories) {
			return nil
		}

		// is text file check
		if fileType == TextFiles {
			isTextFile, err := IsTextFile(filePath)
			if err != nil || !isTextFile {
				return err
			}
		}

		if err := handler(CreateFileSystemFile("",
			CreateFileMetadata(fileInfo.Name(), resolveDirectoryPath(filePath, isDir), filePath, "", fileInfo.ModTime(), size, isDir))); err != nil {
			return err
		}

		return nil
	})
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
