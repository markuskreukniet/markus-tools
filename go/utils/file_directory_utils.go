package utils

import (
	"io"
	"os"
	"path/filepath"
	"time"
	"unicode"
)

// type FileSystemFileExtra struct {
// 	FileSystemFile FileSystemFile
// 	Hash           string
// }

// type FileSystemFile struct {
// 	Data         string
// 	FilePath     string
// 	FileMetadata FileMetadata
// }

// TODO: FileData is bad naming. A hash and metadata is not file data. // FSFile is better naming // File might become a reserved keyword and it might be already used in Go standard libraries
// FileData holds comprehensive information about a file or directory.
// The Identifier field can store either the actual content of the file or a hash of it,
// which makes it useful for various purposes, including as an identifier in unit tests.
type FileData struct {
	Identifier   string // Content or hash of the file
	FileMetadata FileMetadata
}

func CreateFileData(identifier string, file FileMetadata) FileData {
	return FileData{
		Identifier:   identifier,
		FileMetadata: file,
	}
}

type FilesDataGroup struct {
	Identifier string
	FilesData  []FileData
}

func CreateFilesDataGroup(identifier string, files []FileData) FilesDataGroup {
	return FilesDataGroup{
		Identifier: identifier,
		FilesData:  files,
	}
}

type FilesDataGroups []FilesDataGroup

func (groups FilesDataGroups) DidAppendByFileDataIdentifier(file FileData) bool {
	for i, group := range groups {
		if file.Identifier == group.Identifier {
			groups[i].FilesData = append(groups[i].FilesData, file)
			return true
		}
	}
	return false
}

// TODO: For some use cases, FileMetadata has too many fields. Maybe we can use interfaces to solve that problem. Or maybe use viewModels
type FileMetadata struct {
	Path         string
	Name         string
	TimeModified time.Time
	Size         int64
	IsDirectory  bool // It should be a file type, but there is no use case.
}

func (metadata FileMetadata) IsNonZeroByte() bool {
	return metadata.Size > 0
}

func CreateFileMetadata(path, name string, timeModified time.Time, size int64, isDirectory bool) FileMetadata {
	return FileMetadata{
		Path:         path,
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

func ToFileMetadata(filePath string) (FileMetadata, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return FileMetadata{}, err
	}
	return CreateFileMetadata(filePath, info.Name(), info.ModTime(), info.Size(), info.IsDir()), nil
}

func WalkFilterAndHandleFileMetadata(rootFilePath string, mode fileFilterMode, fileType fileType, handler func(FileMetadata)) error {
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
		if size == 0 && (mode == FilesWithoutZeroByteFiles || mode == FilesAndDirectoriesWithoutZeroByteFiles) {
			return nil
		}

		// file type check
		if fileType == PlainTextFiles {
			isTextFile, err := IsNonZeroByteFileATextFile(filePath)
			if err != nil || !isTextFile {
				return err
			}
		}

		handler(CreateFileMetadata(filePath, fileInfo.Name(), fileInfo.ModTime(), size, isDir))
		return nil
	})
}

func FilterAndHandleAllNodesFileMetadata(nodes []FileSystemNode, mode fileFilterMode, handler func(FileMetadata)) error {
	for _, node := range nodes {
		if node.IsDirectory {
			if err := WalkFilterAndHandleFileMetadata(node.Path, mode, AllFiles, func(file FileMetadata) {
				handler(file)
			}); err != nil {
				return err
			}
		} else {
			file, err := ToFileMetadata(node.Path)
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
