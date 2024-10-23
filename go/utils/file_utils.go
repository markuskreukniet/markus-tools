package utils

import (
	"io"
	"os"
	"path/filepath"
	"time"
	"unicode"
)

type FileInfo interface {
	Size() int64
	AbsolutePath() string
}

// MinimalFileInfo implements FileInfo
type MinimalFileInfo struct {
	size         int64
	absolutePath string
}

func (info MinimalFileInfo) Size() int64 {
	return info.size
}

func (info MinimalFileInfo) AbsolutePath() string {
	return info.absolutePath
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

func (info CompleteFileInfo) Size() int64 {
	return info.size
}

func (info CompleteFileInfo) AbsolutePath() string {
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

func toDirectoryPath(filePath string, isDirectory bool) string {
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
		CreateFileMetadata(info.Name(), toDirectoryPath(filePath, isDirectory), filePath, "", info.ModTime(), info.Size(), isDirectory)), nil
}

func FilterAndHandleFileInfo(
	info os.FileInfo, mode fileFilterMode, fileType fileType, absoluteFilePath string, handler func(FileInfo) error,
) error {
	isDir := info.IsDir()

	var size int64
	if !isDir {
		size = info.Size()
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
		isTextFile, err := IsTextFile(absoluteFilePath)
		if err != nil || !isTextFile {
			return err
		}
	}

	handler(CompleteFileInfo{
		name:                  info.Name(),
		absoluteDirectoryPath: toDirectoryPath(absoluteFilePath, isDir),
		absolutePath:          absoluteFilePath,
		timeModified:          info.ModTime(),
		size:                  size,
		isDirectory:           isDir,
	})

	return nil
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
			CreateFileMetadata(fileInfo.Name(), toDirectoryPath(filePath, isDir), filePath, "", fileInfo.ModTime(), size, isDir))); err != nil {
			return err
		}

		return nil
	})
}

func AppendNonZeroByteFiles(nodes []FileSystemNode, files *[]FileSystemFile) error {
	handler := func(file FileSystemFile) error {
		*files = append(*files, file)
		return nil
	}

	for _, node := range nodes {
		if node.IsDirectory {
			if err := WalkFilterAndHandleFileSystemFile(node.Path, NonZeroByteFiles, AllFiles, handler); err != nil {
				return err
			}
		} else {
			file, err := ToFileSystemFile(node.Path)
			if err != nil {
				return err
			}
			handler(file)
		}
	}

	return nil
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
