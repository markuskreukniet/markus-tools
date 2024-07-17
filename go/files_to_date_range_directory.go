package main

import (
	"fmt"
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

func isWithinThreeDays(olderTime, newerTime time.Time) bool {
	return olderTime.Sub(newerTime).Hours() <= 72
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

// Each handler loops unfilteredFiles, but the code is clean enough.
// garbage collection: groups
func filterAndDeleteDuplicateFiles(files []utils.FileData, destinationDirectory string) ([]utils.FileData, error) {
	groups, err := utils.CreateFileHashGroups(files, false)
	if err != nil {
		return nil, err
	}

	files = nil

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
			child, err := isDirectoryChild(destinationDirectory, directory)
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
			child, err := isDirectoryChild(destinationDirectory, file.FileMetadata.Path)
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
			// take the first file and the delete other files
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

func filesToDateRangeDirectory(uniqueFileSystemNodes []utils.FileSystemNode, destinationDirectory string) error {
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

	if err := appendPathsAndFilesByReadingDirectory(destinationDirectory, handler, &files, nil); err != nil {
		return err
	}

	handler = func(_, path string, _ *[]string) {
		badDirectoryFilePaths = append(badDirectoryFilePaths, path)
	}

	for _, path := range goodDirectoryFilePaths {
		appendPathsAndFilesByReadingDirectory(path, handler, &files, nil)
	}

	for _, path := range badDirectoryFilePaths {
		appendPathsAndFilesByReadingDirectoryTree(path, &badDirectoryFilePaths, &files)
	}

	// TODO: is correct?
	if err := utils.AppendNonZeroByteFiles(uniqueFileSystemNodes, &files); err != nil {
		return err
	}

	files, err := filterAndDeleteDuplicateFiles(files, destinationDirectory)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].FileMetadata.TimeModified.Before(files[j].FileMetadata.TimeModified)
	})

	startDateRange := 0
	isFindingDateRange := false
	length := len(files)
	for i := 0; i < length; i++ {
		if i < length-1 && isWithinThreeDays(files[i].FileMetadata.TimeModified, files[i+1].FileMetadata.TimeModified) && !isFindingDateRange {
			isFindingDateRange = true
			startDateRange = i
		} else {
			var name string
			if isFindingDateRange {
				name = createDirectoryDateRangeName(files[startDateRange].FileMetadata.TimeModified, files[i].FileMetadata.TimeModified)
				isFindingDateRange = false
			} else {
				name = toDateFormat(files[i].FileMetadata.TimeModified)
			}

			index := -1
			for j, path := range goodDirectoryFilePaths {
				if strings.HasSuffix(path, name) {
					index = j
					break
				}
			}

			if index == -1 {
				path := filepath.Join(destinationDirectory, name)
				if err := utils.CreateDirectory(path); err != nil {
					return err
				}

				// add files
				for j := startDateRange; j <= i; j++ {
					if err := os.Rename(files[j].FileMetadata.Path, filepath.Join(path, files[j].FileMetadata.Name)); err != nil {
						return err
					}
				}
			} else {
				path := goodDirectoryFilePaths[index]

				// add files
				for j := startDateRange; j <= i; j++ {
					fullPath := filepath.Join(path, files[j].FileMetadata.Name)
					exists, err := utils.FileOrDirectoryExists(fullPath)
					if err != nil {
						return err
					}
					if !exists {
						if err := os.Rename(files[j].FileMetadata.Path, fullPath); err != nil {
							return err
						}
					}
				}

				goodDirectoryFilePaths[index] = goodDirectoryFilePaths[len(goodDirectoryFilePaths)-1]
				goodDirectoryFilePaths = goodDirectoryFilePaths[:len(goodDirectoryFilePaths)-1]
			}
		}
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
