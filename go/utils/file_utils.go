package utils

import (
	"io"
	"os"
	"path/filepath"
	"time"
	"unicode"
)

type FileSystemFileExtra struct {
	Hash           string
	FileSystemFile FileSystemFile
}

func CreateFileSystemFileExtra(hash string, file FileSystemFile) FileSystemFileExtra {
	return FileSystemFileExtra{
		Hash:           hash,
		FileSystemFile: file,
	}
}

type FileSystemFile struct {
	Data         string
	Path         string
	FileMetadata FileMetadata
}

func CreateFileSystemFile(data, filePath string, metadata FileMetadata) FileSystemFile {
	return FileSystemFile{
		Data:         data,
		Path:         filePath,
		FileMetadata: metadata,
	}
}

type FileMetadata struct {
	Name         string
	TimeModified time.Time
	Size         int64
	IsDirectory  bool // It should be a file type, but there is no use case.
}

func CreateFileMetadata(path, name string, timeModified time.Time, size int64, isDirectory bool) FileMetadata {
	return FileMetadata{
		Name:         name,
		TimeModified: timeModified,
		Size:         size,
		IsDirectory:  isDirectory,
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

// TODO: is NonZeroByteFiles and NonZeroByteFilesAndDirectories better naming?
const (
	files fileFilterMode = iota
	FilesWithoutZeroByteFiles
	FilesAndDirectories
	FilesAndDirectoriesWithoutZeroByteFiles
	Directories
)

const (
	AllFiles fileType = iota
	PlainTextFiles
)

const FilePathSeparator = string(filepath.Separator)

func CreateDirectory(filePath string) error {
	return os.Mkdir(filePath, 0755)
}

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

func ToFileSystemFile(filePath string) (FileSystemFile, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return FileSystemFile{}, err
	}

	return CreateFileSystemFile("", filePath, CreateFileMetadata(filePath, info.Name(), info.ModTime(), info.Size(), info.IsDir())), nil
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
		if isDir && (mode == files || mode == FilesWithoutZeroByteFiles) {
			return nil
		}

		// zero byte file check
		if !isDir && size == 0 && (mode == FilesWithoutZeroByteFiles || mode == FilesAndDirectoriesWithoutZeroByteFiles) {
			return nil
		}

		// file type check
		if fileType == PlainTextFiles {
			isTextFile, err := IsNonZeroByteFileATextFile(filePath)
			if err != nil || !isTextFile {
				return err
			}
		}

		if err := handler(CreateFileSystemFile("", filePath, CreateFileMetadata("", fileInfo.Name(), fileInfo.ModTime(), size, isDir))); err != nil {
			return err
		}

		return nil
	})
}

func AppendNonZeroByteFiles(nodes []FileSystemNode, files *[]FileSystemFileExtra) error {
	handler := func(file FileSystemFile) error {
		*files = append(*files, CreateFileSystemFileExtra("", file))
		return nil
	}

	for _, node := range nodes {
		if node.IsDirectory {
			if err := WalkFilterAndHandleFileSystemFile(node.Path, FilesWithoutZeroByteFiles, AllFiles, handler); err != nil {
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

func FileOrDirectoryExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
