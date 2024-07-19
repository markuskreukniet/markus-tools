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

func deleteFiles(files []utils.FileData) error {
	for _, file := range files {
		if err := os.Remove(file.FileMetadata.Path); err != nil {
			return err
		}
	}
	return nil
}

func filterAndDeleteRemainderFiles(files *[]utils.FileData, handler func([]utils.FileData, *[]utils.FileData, *[]utils.FileData) error) error {
	var filteredFiles, remainderFiles []utils.FileData

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
func createFileHandlers(filePath string) []func([]utils.FileData, *[]utils.FileData, *[]utils.FileData) error {
	var handlers []func([]utils.FileData, *[]utils.FileData, *[]utils.FileData) error

	// shortest file name
	handler := func(unfilteredFiles []utils.FileData, filteredFiles, remainderFiles *[]utils.FileData) error {
		minimumLength := 0

		for _, file := range unfilteredFiles {
			length := len(file.FileMetadata.Name)
			if length < minimumLength || minimumLength == 0 {
				minimumLength = length
				*remainderFiles = append(*remainderFiles, *filteredFiles...)
				*filteredFiles = []utils.FileData{file}
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
	handler = func(unfilteredFiles []utils.FileData, filteredFiles, remainderFiles *[]utils.FileData) error {
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
	handler = func(unfilteredFiles []utils.FileData, filteredFiles, remainderFiles *[]utils.FileData) error {
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
	handler = func(unfilteredFiles []utils.FileData, filteredFiles, remainderFiles *[]utils.FileData) error {
		var newestTime time.Time

		for _, file := range unfilteredFiles {
			if file.FileMetadata.TimeModified.After(newestTime) {
				newestTime = file.FileMetadata.TimeModified
				*remainderFiles = append(*remainderFiles, *filteredFiles...)
				*filteredFiles = []utils.FileData{file}
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
func filterAndDeleteDuplicateFiles(files []utils.FileData, destinationDirectory string) ([]utils.FileData, error) {
	groups, err := utils.CreateFileHashGroups(files, false)
	if err != nil {
		return nil, err
	}

	files = nil

	handlers := createFileHandlers(destinationDirectory)

	for _, group := range groups {
		for _, handler := range handlers {
			if len(group.FilesData) > 1 {
				err := filterAndDeleteRemainderFiles(&group.FilesData, handler)
				if err != nil {
					return nil, err
				}
			} else {
				files = append(files, group.FilesData[0])
				break
			}
		}
		if len(group.FilesData) > 1 {
			// append first file and the delete other files
			files = append(files, group.FilesData[0])
			group.FilesData[0] = group.FilesData[len(group.FilesData)-1]
			group.FilesData = group.FilesData[:len(group.FilesData)-1]
			if err := deleteFiles(group.FilesData); err != nil {
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

func appendPathsAndFilesByReadingDirectoryTree(path string, paths *[]string, files *[]utils.FileData) error {
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

func appendPathsAndFilesByReadingDirectory(path string, handler func(string, string, *[]string), files *[]utils.FileData, stack *[]string) error {
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
			*files = append(*files, utils.CreateFileData("", utils.CreateFileMetadata(fullPath, info.Name(), info.ModTime(), info.Size(), false)))
		}
	}
	return nil
}

// garbage collection: handler
func createFilesAndDirectoryFilePaths(filePath string) ([]utils.FileData, []string, []string, error) {
	var files []utils.FileData
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

// TODO: cleaning
// garbage collection: startDateRange, isFindingDateRange
func moveFilesToDateRangeDirectoriesAndRemoveUsedGoodDirectories(files []utils.FileData, filePaths []string, filePath string) ([]string, error) {
	length := len(files)
	startDateRange := 0
	isFindingDateRange := true

	if length == 0 {
		return filePaths, nil
	}

	for i := 1; i < length; i++ {
		iMinusOne := i - 1
		if isWithin72Hours(files[iMinusOne].FileMetadata.TimeModified, files[i].FileMetadata.TimeModified) {
			if !isFindingDateRange {
				isFindingDateRange = true
				startDateRange = iMinusOne
			}
			continue
		}

		var name string
		if isFindingDateRange {
			name = createDirectoryDateRangeName(files[startDateRange].FileMetadata.TimeModified, files[i].FileMetadata.TimeModified)
			isFindingDateRange = false
		} else {
			name = toDateFormat(files[i].FileMetadata.TimeModified)
			startDateRange = i
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
			utils.CreateDirectory(directoryFilePath)
		}

		// add files
		for j := startDateRange; j <= i; j++ {
			fullFilePath := filepath.Join(directoryFilePath, files[j].FileMetadata.Name)
			exists, err := utils.FileOrDirectoryExists(fullFilePath)
			if err != nil {
				return nil, err
			}
			if !exists {
				if err := os.Rename(files[j].FileMetadata.Path, fullFilePath); err != nil {
					return nil, err
				}
			}
		}
	}

	if isFindingDateRange {
		lengthMinusOne := length - 1

		name := createDirectoryDateRangeName(files[startDateRange].FileMetadata.TimeModified, files[lengthMinusOne].FileMetadata.TimeModified)

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
			utils.CreateDirectory(directoryFilePath)
		}

		// add files
		for j := startDateRange; j <= lengthMinusOne; j++ {
			fullFilePath := filepath.Join(directoryFilePath, files[j].FileMetadata.Name)
			exists, err := utils.FileOrDirectoryExists(fullFilePath)
			if err != nil {
				return nil, err
			}
			if !exists {
				if err := os.Rename(files[j].FileMetadata.Path, fullFilePath); err != nil {
					return nil, err
				}
			}
		}
	}

	return filePaths, nil
}

// TODO: renaming + cleaning
// garbage collection: startDateRange, isFindingDateRange, length
func moveFilesToDateRangeDirectoriesAndFilterDirectories(files []utils.FileData, filePaths []string, filePath string) ([]string, error) {
	startDateRange := 0
	isFindingDateRange := false
	length := len(files)

	for i := 0; i < length; i++ {
		if i < length-1 && isWithin72Hours(files[i].FileMetadata.TimeModified, files[i+1].FileMetadata.TimeModified) && !isFindingDateRange {
			isFindingDateRange = true
			startDateRange = i
		} else {
			// directory name
			//TODO: wrong comment
			var name string
			if isFindingDateRange {
				name = createDirectoryDateRangeName(files[startDateRange].FileMetadata.TimeModified, files[i].FileMetadata.TimeModified)
				isFindingDateRange = false
			} else {
				startDateRange = i
				name = toDateFormat(files[i].FileMetadata.TimeModified)
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

			// TODO: should CreateDirectory create a dir with the same rights?
			if !isDirectoryFound {
				utils.CreateDirectory(directoryFilePath)
			}

			// add files
			for j := startDateRange; j <= i; j++ {
				fullFilePath := filepath.Join(directoryFilePath, files[j].FileMetadata.Name)
				exists, err := utils.FileOrDirectoryExists(fullFilePath)
				if err != nil {
					return nil, err
				}
				if !exists {
					if err := os.Rename(files[j].FileMetadata.Path, fullFilePath); err != nil {
						return nil, err
					}
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

	// TODO: is correct?
	if err := utils.AppendNonZeroByteFiles(uniqueFileSystemNodes, &files); err != nil {
		return err
	}

	files, err = filterAndDeleteDuplicateFiles(files, destinationDirectory)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].FileMetadata.TimeModified.Before(files[j].FileMetadata.TimeModified)
	})

	// TODO: goodDirectoryFilePaths should work with reference?
	goodDirectoryFilePaths, err = moveFilesToDateRangeDirectoriesAndFilterDirectories(files, goodDirectoryFilePaths, destinationDirectory)
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
