package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

const spacedHyphen = " - "
const dateLayout = "2006-01-02" // YYYY-MM-DD

func isValidDateRangeDirectoryName(name string) bool {
	if strings.Contains(name, spacedHyphen) {
		nameParts := strings.Split(name, spacedHyphen)
		firstDate, err := parseDate(nameParts[0])
		if err != nil {
			return false
		}
		secondDate, err := parseDate(nameParts[1])
		if err != nil {
			return false
		}
		daysDifference := secondDate.Sub(firstDate).Hours() / 24
		if daysDifference >= 1 {
			return true
		}
	} else if isValidDateFormat(name) {
		return true
	}
	return false
}

func isWithin72Hours(olderTime, newerTime time.Time) bool {
	return math.Abs(olderTime.Sub(newerTime).Hours()) <= 72
}

func createDirectoryDateRangeName(startTime, endTime time.Time) string {
	start := toDateFormat(startTime)
	end := toDateFormat(endTime)

	if start == end {
		return start
	}
	return fmt.Sprintf("%s - %s", start, end)
}

func deleteFiles(files []utils.FileSystemFile) error {
	for _, file := range files {
		if err := os.Remove(file.FileMetadata.Path); err != nil {
			return err
		}
	}
	return nil
}

func filterAndDeleteRemainderFiles(files *[]utils.FileSystemFile, handler func([]utils.FileSystemFile, *[]utils.FileSystemFile, *[]utils.FileSystemFile) error) error {
	var filteredFiles, remainderFiles []utils.FileSystemFile

	err := handler(*files, &filteredFiles, &remainderFiles)
	if err != nil {
		return err
	}

	if len(filteredFiles) > 0 {
		*files = filteredFiles
		if err := deleteFiles(remainderFiles); err != nil {
			return err
		}
	}

	return nil
}

// garbage collection: handler
func createFileHandlers(filePath string) []func([]utils.FileSystemFile, *[]utils.FileSystemFile, *[]utils.FileSystemFile) error {
	var handlers []func([]utils.FileSystemFile, *[]utils.FileSystemFile, *[]utils.FileSystemFile) error

	// shortest file name
	handler := func(unfilteredFiles []utils.FileSystemFile, filteredFiles, remainderFiles *[]utils.FileSystemFile) error {
		minimumLength := 0

		for _, file := range unfilteredFiles {
			length := len(file.FileMetadata.Name)
			if length < minimumLength || minimumLength == 0 {
				minimumLength = length
				*remainderFiles = append(*remainderFiles, *filteredFiles...)
				*filteredFiles = []utils.FileSystemFile{file}
			} else if length == minimumLength {
				*filteredFiles = append(*filteredFiles, file)
			} else {
				*remainderFiles = append(*remainderFiles, file)
			}
		}

		return nil
	}
	handlers = append(handlers, handler)

	// valid name of date directory or date range directory
	handler = func(unfilteredFiles []utils.FileSystemFile, filteredFiles, remainderFiles *[]utils.FileSystemFile) error {
		for _, file := range unfilteredFiles {
			directory := filepath.Dir(file.FileMetadata.Path)
			child, err := isDirectoryChild(filePath, directory)
			if err != nil {
				return err
			}
			if child && isValidDateRangeDirectoryName(directory) {
				*filteredFiles = append(*filteredFiles, file)
			} else {
				*remainderFiles = append(*remainderFiles, file)
			}
		}

		return nil
	}
	handlers = append(handlers, handler)

	// destination directory
	handler = func(unfilteredFiles []utils.FileSystemFile, filteredFiles, remainderFiles *[]utils.FileSystemFile) error {
		for _, file := range unfilteredFiles {
			child, err := isDirectoryChild(filePath, file.FileMetadata.Path)
			if err != nil {
				return err
			}
			if child {
				*filteredFiles = append(*filteredFiles, file)
			} else {
				*remainderFiles = append(*remainderFiles, file)
			}
		}

		return nil
	}
	handlers = append(handlers, handler)

	// newest modification time
	handler = func(unfilteredFiles []utils.FileSystemFile, filteredFiles, remainderFiles *[]utils.FileSystemFile) error {
		var newestTime time.Time

		for _, file := range unfilteredFiles {
			if file.FileMetadata.TimeModified.After(newestTime) {
				newestTime = file.FileMetadata.TimeModified
				*remainderFiles = append(*remainderFiles, *filteredFiles...)
				*filteredFiles = []utils.FileSystemFile{file}
			} else if file.FileMetadata.TimeModified.Equal(newestTime) {
				*filteredFiles = append(*filteredFiles, file)
			} else {
				*remainderFiles = append(*remainderFiles, file)
			}
		}

		return nil
	}
	handlers = append(handlers, handler)

	return handlers
}

// Each handler loops unfilteredFiles, but the code is clean enough.
// garbage collection: groups
func filterAndDeleteDuplicateFiles(files []utils.FileSystemFile, destinationDirectory string) ([]utils.FileSystemFile, error) {
	groups, err := utils.CreateFileSystemFileByHashGroups(files, false)
	if err != nil {
		return nil, err
	}

	files = nil

	handlers := createFileHandlers(destinationDirectory)

	for _, group := range groups {
		for _, handler := range handlers {
			if len(group) > 1 {
				err := filterAndDeleteRemainderFiles(&group, handler)
				if err != nil {
					return nil, err
				}
			} else {
				files = append(files, group[0])
				break
			}
		}
		if len(group) > 1 {
			// append first file and the delete other files
			files = append(files, group[0])
			group[0] = group[len(group)-1]
			group = group[:len(group)-1]
			if err := deleteFiles(group); err != nil {
				return nil, err
			}
		}
	}

	return files, nil
}

func isDirectoryChild(filePath, childFilePath string) (bool, error) {
	path, err := filepath.Rel(filePath, childFilePath)
	if err != nil {
		return false, err
	}
	return !strings.HasPrefix(path, "..") && !strings.Contains(path, utils.FilePathSeparator), nil
}

// TODO: naming
type dateRangeArg struct {
	directoryName string
	filePath      string
}

func addDirectory(directories *[]string, arg dateRangeArg) {
	*directories = append(*directories, arg.filePath)
}

func categorizeFilesAndDirectories(destinationDirectory string) ([]utils.DateRangeFileInfo, []string, []string, error) {
	var files []utils.DateRangeFileInfo
	var goodDirectoryPaths []string
	var badDirectoryPaths []string

	categorizeInDirectory := func(directoryPaths *[]string, arg dateRangeArg) {
		if isValidDateRangeDirectoryName(arg.directoryName) {
			goodDirectoryPaths = append(goodDirectoryPaths, arg.filePath)
		} else {
			*directoryPaths = append(*directoryPaths, arg.filePath)
		}
	}

	entries, err := os.ReadDir(destinationDirectory)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, nil, nil, err
		}
		categorize(
			info, filepath.Join(destinationDirectory, entry.Name()), &files, &badDirectoryPaths, categorizeInDirectory,
		)
	}

	for _, path := range append(goodDirectoryPaths, badDirectoryPaths...) {
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if path != filePath {
				categorize(info, filePath, &files, &badDirectoryPaths, addDirectory)
			}

			return nil
		})
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return files, goodDirectoryPaths, badDirectoryPaths, nil
}

func categorize(
	info os.FileInfo,
	filePath string,
	files *[]utils.DateRangeFileInfo,
	badDirectoryPaths *[]string,
	handler func(*[]string, dateRangeArg),
) error {
	if info.IsDir() {
		handler(badDirectoryPaths, dateRangeArg{
			directoryName: filepath.Base(filePath),
			filePath:      filePath,
		})
	} else if info.Mode().IsRegular() {
		size := info.Size()
		if size > 0 {
			*files = append(*files, utils.DateRangeFileInfo{
				Size:         size,
				Path:         filePath,
				TimeModified: info.ModTime(),
			})
		} else {
			// TODO: error
		}
	} else {
		// TODO: error
	}

	return nil
}

// TODO: Does not work efficient, could be done without making groups?
// garbage collection: length, groups, groupIndex
func moveFilesToDateRangeDirectoriesAndRemoveUsedGoodDirectories(files []utils.FileSystemFile, filePaths []string, filePath string) ([]string, error) {
	length := len(files)

	if length == 0 {
		return filePaths, nil
	}

	groups := [][]utils.FileSystemFile{{files[0]}}
	groupIndex := 0

	for i := 1; i < length; i++ {
		iMinusOne := i - 1
		if isWithin72Hours(files[iMinusOne].FileMetadata.TimeModified, files[i].FileMetadata.TimeModified) {
			groups[groupIndex] = append(groups[groupIndex], files[i])
		} else {
			groupIndex++
			groups = append(groups, []utils.FileSystemFile{files[i]})
		}
	}

	for _, group := range groups {
		length = len(group)
		lengthMinusOne := length - 1
		var name string
		if group[0].FileMetadata.TimeModified == group[lengthMinusOne].FileMetadata.TimeModified {
			name = toDateFormat(group[0].FileMetadata.TimeModified)
		} else {
			name = createDirectoryDateRangeName(group[0].FileMetadata.TimeModified, group[lengthMinusOne].FileMetadata.TimeModified)
		}
		directoryFilePath := filepath.Join(filePath, name)
		isDirectoryFound := false
		for j, path := range filePaths {
			if path == directoryFilePath {
				isDirectoryFound = true
				filePaths[j] = filePaths[len(filePaths)-1]
				filePaths = filePaths[:len(filePaths)-1]
				break
			}
		}

		// TODO: should CreateDirectory create a dir with the same rights as parent dir?
		if !isDirectoryFound {
			if err := utils.CreateDirectory(directoryFilePath); err != nil {
				return nil, err
			}
		}

		// TODO clean and make it more efficient
		// add files
		for _, file := range group {
			fullFilePath := filepath.Join(directoryFilePath, file.FileMetadata.Name)
			exists, err := utils.FileExists(fullFilePath)
			if err != nil {
				return nil, err
			}
			if exists {
				// We should always create a hash of the file in the destination folder.
				// Otherwise, we have to loop through all the files to find that file, and that found file might not have a hash yet.
				hash, err := utils.CreateFileHash(fullFilePath)
				if err != nil {
					return nil, err
				}
				if file.FileMetadata.Hash == "" {
					file.FileMetadata.Hash, err = utils.CreateFileHash(file.FileMetadata.Path)
					if err != nil {
						return nil, err
					}
				}
				if hash != file.FileMetadata.Hash {
					extension := filepath.Ext(file.FileMetadata.Name)
					nameWithoutExtension := strings.TrimSuffix(file.FileMetadata.Name, extension)
					fullFilePath = filepath.Join(directoryFilePath, nameWithoutExtension+" 2"+extension)

					if err := os.Rename(file.FileMetadata.Path, fullFilePath); err != nil {
						return nil, err
					}
				}
			} else {
				if err := os.Rename(file.FileMetadata.Path, fullFilePath); err != nil {
					return nil, err
				}
			}
		}
	}

	return filePaths, nil
}

func createHandlers(
	destinationDirectory string) []func([]utils.DateRangeFileInfo, *[]utils.DateRangeFileInfo,
) []utils.DateRangeFileInfo {
	appendBadFilesAndReplaceGoodFiles := func(
		badFiles *[]utils.DateRangeFileInfo, goodFiles *[]utils.DateRangeFileInfo, file utils.DateRangeFileInfo,
	) {
		*badFiles = append(*badFiles, *goodFiles...)
		goodFiles = &[]utils.DateRangeFileInfo{file}
	}

	categorizeOnShortestFileNameLength := func(
		files []utils.DateRangeFileInfo, badFiles *[]utils.DateRangeFileInfo,
	) []utils.DateRangeFileInfo {
		getNameLength := func(file utils.DateRangeFileInfo) int {
			return len(filepath.Base(file.Path))
		}

		good := []utils.DateRangeFileInfo{files[0]}
		var minimumLength = getNameLength(files[0])

		for i := 1; i < len(files); i++ {
			nameLength := getNameLength(files[i])
			if nameLength < minimumLength {
				minimumLength = nameLength
				appendBadFilesAndReplaceGoodFiles(badFiles, &good, files[i])
			} else if nameLength == minimumLength {
				good = append(good, files[i])
			} else {
				*badFiles = append(*badFiles, files[i])
			}
		}

		return good
	}

	categorizeOnValidDateRangeDirectoryName := func(
		files []utils.DateRangeFileInfo, badFiles *[]utils.DateRangeFileInfo,
	) []utils.DateRangeFileInfo {
		var tempGood1Files []utils.DateRangeFileInfo
		var tempGood2Files []utils.DateRangeFileInfo
		var tempBadFiles []utils.DateRangeFileInfo

		for _, file := range files {
			directoryPath := filepath.Dir(file.Path)
			if filepath.Dir(directoryPath) == destinationDirectory {
				if isValidDateRangeDirectoryName(directoryPath) {
					tempGood2Files = append(tempGood2Files, file)
				} else {
					tempGood1Files = append(tempGood1Files, file)
				}
			} else {
				tempBadFiles = append(tempBadFiles, file)
			}
		}

		if len(tempGood2Files) > 0 {
			*badFiles = append(*badFiles, tempGood1Files...)
			*badFiles = append(*badFiles, tempBadFiles...)
			return tempGood2Files
		}

		if len(tempGood1Files) > 0 {
			*badFiles = append(*badFiles, tempBadFiles...)
			return tempGood1Files
		}

		return tempBadFiles
	}

	return []func([]utils.DateRangeFileInfo, *[]utils.DateRangeFileInfo) []utils.DateRangeFileInfo{
		categorizeOnShortestFileNameLength, categorizeOnValidDateRangeDirectoryName,
	}
}

// func deleteDuplicateFiles(files []utils.DateRangeFileInfo, destinationDirectory string) error {
// 	groups, err := utils.CreateDuplicateFileInfoGroupsByHash(files, false)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func filesToDateRangeDirectory(uniqueFileSystemNodes []utils.FileSystemNode, destinationDirectory string) error {
	filesNew, goodDirectoryFilePaths, badDirectoryFilePaths, err := categorizeFilesAndDirectories(destinationDirectory)
	if err != nil {
		return err
	}

	for _, node := range uniqueFileSystemNodes {
		info, err := os.Stat(node.Path)
		if err != nil {
			return err
		}
		categorize(info, node.Path, &filesNew, &badDirectoryFilePaths, addDirectory)
	}

	// TODO: remove this converting
	var files []utils.FileSystemFile
	for _, file := range filesNew {
		files = append(files, utils.FileSystemFile{
			Data:         "",
			FileMetadata: utils.CreateFileMetadata(filepath.Base(file.Path), filepath.Dir(file.Path), file.Path, "", file.TimeModified, file.Size, false),
		})
	}

	files, err = filterAndDeleteDuplicateFiles(files, destinationDirectory)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].FileMetadata.TimeModified.Before(files[j].FileMetadata.TimeModified)
	})

	// TODO: goodDirectoryFilePaths should work with reference?
	goodDirectoryFilePaths, err = moveFilesToDateRangeDirectoriesAndRemoveUsedGoodDirectories(files, goodDirectoryFilePaths, destinationDirectory)
	if err != nil {
		return err
	}

	// Remove the bad empty directories
	// There is no need to check if the directory exists before attempting removal.
	for i := len(badDirectoryFilePaths) - 1; i >= 0; i-- {
		if err := os.Remove(badDirectoryFilePaths[i]); err != nil {
			return err
		}
	}

	for _, path := range goodDirectoryFilePaths {
		if err := os.Remove(path); err != nil {
			return err
		}
	}

	return nil
}

func toDateFormat(time time.Time) string {
	return time.Format(dateLayout)
}

func parseDate(rawDate string) (time.Time, error) {
	date, err := time.Parse(dateLayout, rawDate)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

func isValidDateFormat(rawDate string) bool {
	_, err := parseDate(rawDate)
	return err == nil
}
