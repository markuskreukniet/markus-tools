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
		if daysDifference >= 1 && daysDifference <= 3 {
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

func deleteFiles(files []utils.FileSystemFileExtra) error {
	for _, file := range files {
		if err := os.Remove(file.FileSystemFile.Path); err != nil {
			return err
		}
	}
	return nil
}

func filterAndDeleteRemainderFiles(files *[]utils.FileSystemFileExtra, handler func([]utils.FileSystemFileExtra, *[]utils.FileSystemFileExtra, *[]utils.FileSystemFileExtra) error) error {
	var filteredFiles, remainderFiles []utils.FileSystemFileExtra

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
func createFileHandlers(filePath string) []func([]utils.FileSystemFileExtra, *[]utils.FileSystemFileExtra, *[]utils.FileSystemFileExtra) error {
	var handlers []func([]utils.FileSystemFileExtra, *[]utils.FileSystemFileExtra, *[]utils.FileSystemFileExtra) error

	// shortest file name
	handler := func(unfilteredFiles []utils.FileSystemFileExtra, filteredFiles, remainderFiles *[]utils.FileSystemFileExtra) error {
		minimumLength := 0

		for _, file := range unfilteredFiles {
			length := len(file.FileSystemFile.FileMetadata.Name)
			if length < minimumLength || minimumLength == 0 {
				minimumLength = length
				*remainderFiles = append(*remainderFiles, *filteredFiles...)
				*filteredFiles = []utils.FileSystemFileExtra{file}
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
	handler = func(unfilteredFiles []utils.FileSystemFileExtra, filteredFiles, remainderFiles *[]utils.FileSystemFileExtra) error {
		for _, file := range unfilteredFiles {
			directory := filepath.Dir(file.FileSystemFile.Path)
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
	handler = func(unfilteredFiles []utils.FileSystemFileExtra, filteredFiles, remainderFiles *[]utils.FileSystemFileExtra) error {
		for _, file := range unfilteredFiles {
			child, err := isDirectoryChild(filePath, file.FileSystemFile.Path)
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
	handler = func(unfilteredFiles []utils.FileSystemFileExtra, filteredFiles, remainderFiles *[]utils.FileSystemFileExtra) error {
		var newestTime time.Time

		for _, file := range unfilteredFiles {
			if file.FileSystemFile.FileMetadata.TimeModified.After(newestTime) {
				newestTime = file.FileSystemFile.FileMetadata.TimeModified
				*remainderFiles = append(*remainderFiles, *filteredFiles...)
				*filteredFiles = []utils.FileSystemFileExtra{file}
			} else if file.FileSystemFile.FileMetadata.TimeModified.Equal(newestTime) {
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
func filterAndDeleteDuplicateFiles(files []utils.FileSystemFileExtra, destinationDirectory string) ([]utils.FileSystemFileExtra, error) {
	groups, err := utils.CreateFileSystemFileExtraByHashGroups(files, false)
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

func appendPathsAndFilesByReadingDirectoryTree(path string, paths *[]string, files *[]utils.FileSystemFileExtra) error {
	handler := func(_, path string, stack *[]string) {
		*paths = append(*paths, path)
		*stack = append(*stack, path)
	}

	stack := []string{path}
	for len(stack) > 0 {
		path := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if err := appendPathsAndFilesByReadingDirectory(path, handler, files, &stack); err != nil {
			return err
		}
	}
	return nil
}

func appendPathsAndFilesByReadingDirectory(path string, handler func(string, string, *[]string), files *[]utils.FileSystemFileExtra, stack *[]string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		name := entry.Name()
		fullPath := filepath.Join(path, name)
		if entry.IsDir() {
			handler(name, fullPath, stack)
		} else {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			*files = append(*files,
				utils.CreateFileSystemFileExtra("",
					utils.CreateFileSystemFile("", fullPath,
						utils.CreateFileMetadata("", info.Name(), info.ModTime(), info.Size(), false))))
		}
	}
	return nil
}

// garbage collection: handler
func createFilesAndDirectoryFilePaths(filePath string) ([]utils.FileSystemFileExtra, []string, []string, error) {
	var files []utils.FileSystemFileExtra
	var goodDirectoryFilePaths []string
	var badDirectoryFilePaths []string

	handler := func(name, path string, _ *[]string) {
		if isValidDateRangeDirectoryName(name) {
			goodDirectoryFilePaths = append(goodDirectoryFilePaths, path)
		} else {
			badDirectoryFilePaths = append(badDirectoryFilePaths, path)
		}
	}

	if err := appendPathsAndFilesByReadingDirectory(filePath, handler, &files, nil); err != nil {
		return nil, nil, nil, err
	}

	handler = func(_, path string, _ *[]string) {
		badDirectoryFilePaths = append(badDirectoryFilePaths, path)
	}

	for _, path := range goodDirectoryFilePaths {
		if err := appendPathsAndFilesByReadingDirectory(path, handler, &files, nil); err != nil {
			return nil, nil, nil, err
		}
	}

	for _, path := range badDirectoryFilePaths {
		appendPathsAndFilesByReadingDirectoryTree(path, &badDirectoryFilePaths, &files)
	}

	return files, goodDirectoryFilePaths, badDirectoryFilePaths, nil
}

// TODO: Does not work efficient, could be done without making groups?
// garbage collection: length, groups, groupIndex
func moveFilesToDateRangeDirectoriesAndRemoveUsedGoodDirectories(files []utils.FileSystemFileExtra, filePaths []string, filePath string) ([]string, error) {
	length := len(files)

	if length == 0 {
		return filePaths, nil
	}

	groups := [][]utils.FileSystemFileExtra{{files[0]}}
	groupIndex := 0

	for i := 1; i < length; i++ {
		iMinusOne := i - 1
		if isWithin72Hours(files[iMinusOne].FileSystemFile.FileMetadata.TimeModified, files[i].FileSystemFile.FileMetadata.TimeModified) {
			groups[groupIndex] = append(groups[groupIndex], files[i])
		} else {
			groupIndex++
			groups = append(groups, []utils.FileSystemFileExtra{files[i]})
		}
	}

	for _, group := range groups {
		length = len(group)
		lengthMinusOne := length - 1
		var name string
		if group[0].FileSystemFile.FileMetadata.TimeModified == group[lengthMinusOne].FileSystemFile.FileMetadata.TimeModified {
			name = toDateFormat(group[0].FileSystemFile.FileMetadata.TimeModified)
		} else {
			name = createDirectoryDateRangeName(group[0].FileSystemFile.FileMetadata.TimeModified, group[lengthMinusOne].FileSystemFile.FileMetadata.TimeModified)
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
			fullFilePath := filepath.Join(directoryFilePath, file.FileSystemFile.FileMetadata.Name)
			exists, err := utils.FileOrDirectoryExists(fullFilePath)
			if err != nil {
				return nil, err
			}
			if exists {
				// We should always create a hash of the file in the destination folder.
				// Otherwise, we have to loop through all the files to find that file, and that found file might not have a hash yet.
				hash, err := utils.HashFile(fullFilePath)
				if err != nil {
					return nil, err
				}
				if file.Hash == "" {
					file.Hash, err = utils.HashFile(file.FileSystemFile.Path)
					if err != nil {
						return nil, err
					}
				}
				if hash != file.Hash {
					extension := filepath.Ext(file.FileSystemFile.FileMetadata.Name)
					nameWithoutExtension := strings.TrimSuffix(file.FileSystemFile.FileMetadata.Name, extension)
					fullFilePath = filepath.Join(directoryFilePath, nameWithoutExtension+" 2"+extension)

					if err := os.Rename(file.FileSystemFile.Path, fullFilePath); err != nil {
						return nil, err
					}
				}
			} else {
				if err := os.Rename(file.FileSystemFile.Path, fullFilePath); err != nil {
					return nil, err
				}
			}
		}
	}

	return filePaths, nil
}

func filesToDateRangeDirectory(uniqueFileSystemNodes []utils.FileSystemNode, destinationDirectory string) error {
	files, goodDirectoryFilePaths, badDirectoryFilePaths, err := createFilesAndDirectoryFilePaths(destinationDirectory)
	if err != nil {
		return err
	}

	if err := utils.AppendNonZeroByteFilesNew(uniqueFileSystemNodes, &files); err != nil {
		return err
	}

	files, err = filterAndDeleteDuplicateFiles(files, destinationDirectory)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].FileSystemFile.FileMetadata.TimeModified.Before(files[j].FileSystemFile.FileMetadata.TimeModified)
	})

	// TODO: goodDirectoryFilePaths should work with reference?
	goodDirectoryFilePaths, err = moveFilesToDateRangeDirectoriesAndRemoveUsedGoodDirectories(files, goodDirectoryFilePaths, destinationDirectory)
	if err != nil {
		return err
	}

	// Remove the bad empty directories
	// There is no need to check if the file path exists before attempting removal.
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
