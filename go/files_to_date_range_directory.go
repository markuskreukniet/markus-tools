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

type filePathTimeModified struct {
	filePath     string
	timeModified time.Time
}

const spacedHyphen = " - "
const dateLayout = "2006-01-02" // YYYY-MM-DD

func isValidDateRangeDirectoryName(name string) bool {
	if strings.Contains(name, spacedHyphen) {
		baseParts := strings.Split(name, spacedHyphen)
		firstDate, err := parseDate(baseParts[0])
		if err != nil {
			return false
		}
		secondDate, err := parseDate(baseParts[1])
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

func isValidDateRangeDirectory(filePath string) bool {
	return isValidDateRangeDirectoryName(filepath.Base(filePath))
}

func isWithinThreeDays(olderTime, newerTime time.Time) bool {
	return olderTime.Sub(newerTime).Hours() <= 72
}

func createDirectoryDateRangeName(startTime, endTime time.Time) (string, error) {
	var builder strings.Builder
	if err := formatDateAndWriteString(&builder, startTime); err != nil {
		return "", err
	}
	if _, err := builder.WriteString(spacedHyphen); err != nil {
		return "", err
	}
	if err := formatDateAndWriteString(&builder, endTime); err != nil {
		return "", err
	}
	return builder.String(), nil
}

// garbage collection: groups
func deleteDuplicateFiles(files []utils.FileData) ([]utils.FileData, error) {
	groups, err := utils.CreateDuplicateFileGroups(files)
	if err != nil {
		return nil, err
	}
	files = nil
	for _, group := range groups {
		if len(group.FilesData) > 1 {
			//

			// take the first file
			files = append(files, group.FilesData[0])
		}
	}

	// TODO: delete duplicate files

	return files, nil
}

// TODO: WIP
func filesToDateRangeDirectoryWIP(uniqueFileSystemNodes []utils.FileSystemNode, destinationDirectory string) error {
	var files []utils.FileData
	var goodDirectoryFilePaths []string
	var badDirectoryFilePaths []string

	entries, err := os.ReadDir(destinationDirectory)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(destinationDirectory, entry.Name())
		isDir := entry.IsDir()
		if isDir {
			if isValidDateRangeDirectoryName(entry.Name()) {
				goodDirectoryFilePaths = append(goodDirectoryFilePaths, path)
			} else {
				badDirectoryFilePaths = append(badDirectoryFilePaths, path)
			}
		} else {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			files = append(files, utils.CreateFileData("", utils.CreateFileMetadata(path, info.ModTime(), info.Size(), isDir)))
		}
	}

	handler := func(metadata utils.FileMetadata) {
		if metadata.IsDirectory {
			badDirectoryFilePaths = append(badDirectoryFilePaths, metadata.FilePath)
		} else {
			files = append(files, utils.CreateFileData("", metadata))
		}
	}

	for _, path := range goodDirectoryFilePaths {
		if err := utils.WalkFilterAndHandleFileMetadata(path, utils.FilesAndDirectories, utils.AllFiles, handler); err != nil {
			return err
		}
	}
	for _, path := range badDirectoryFilePaths {
		if err := utils.WalkFilterAndHandleFileMetadata(path, utils.FilesAndDirectories, utils.AllFiles, handler); err != nil {
			return err
		}
	}

	if files, err = deleteDuplicateFiles(files); err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].FileMetadata.ModificationTime.Before(files[j].FileMetadata.ModificationTime)
	})

	startDateRange := 0
	isFindingDateRange := false
	length := len(files)
	for i := 0; i < length; i++ {
		if i < length-1 && isWithinThreeDays(files[i].FileMetadata.ModificationTime, files[i+1].FileMetadata.ModificationTime) && !isFindingDateRange {
			isFindingDateRange = true
			startDateRange = i
		} else {
			var name string
			if isFindingDateRange {
				// Declare 'err' separately to avoid shadowing 'name' with ':='
				var err error
				name, err = createDirectoryDateRangeName(files[startDateRange].FileMetadata.ModificationTime, files[i].FileMetadata.ModificationTime)
				if err != nil {
					return err
				}

				isFindingDateRange = false
			} else {
				name = toDateFormat(files[i].FileMetadata.ModificationTime)
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

				// TODO: is this useful?
				goodDirectoryFilePaths = append(goodDirectoryFilePaths, path)
				if err := utils.CreateDirectory(path); err != nil {
					return err
				}

				// add files
				for j := startDateRange; j <= i; j++ {
					if err := os.Rename(files[j].FileMetadata.FilePath, extractBaseAndJoinWithFilePath(files[j].FileMetadata.FilePath, path)); err != nil {
						return err
					}
				}
			} else {
				path := goodDirectoryFilePaths[index]

				// add files
				for j := startDateRange; j <= i; j++ {
					fullPath := extractBaseAndJoinWithFilePath(files[j].FileMetadata.FilePath, path)
					exists, err := utils.FileOrDirectoryExists(fullPath)
					if err != nil {
						return err
					}
					if !exists {
						if err := os.Rename(files[j].FileMetadata.FilePath, fullPath); err != nil {
							return err
						}
					}
				}
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

	return nil
}

func extractBaseAndJoinWithFilePath(filePathWithBase, filePath string) string {
	return filepath.Join(filePath, filepath.Base(filePathWithBase))
}

func filesToDateRangeDirectory(uniqueFileSystemNodes []utils.FileSystemNode, destinationDirectory string) error {
	var filePathsTimeModified []filePathTimeModified
	if err := appendFilePathsTimeModified(&filePathsTimeModified, uniqueFileSystemNodes); err != nil {
		return err
	}
	var dateRangeDirectoryPaths []string

	// Two error variables are needed because errJ might become an error while errI does not.
	var errJ error
	errI := utils.AppendFileDetails(
		func(detail utils.FileDetail) {
			if isValidDateRangeDirectory(detail.Path) {
				dateRangeDirectoryPaths = append(dateRangeDirectoryPaths, detail.Path)
				errJ = appendFilePathsTimeModified(&filePathsTimeModified, createDirectoryFileSystemNodeInSlice(detail.Path))
			}
		}, createDirectoryFileSystemNodeInSlice(destinationDirectory), utils.Directories)
	if errI != nil {
		return errI
	}
	if errJ != nil {
		return errJ
	}

	sort.Slice(filePathsTimeModified, func(i, j int) bool {
		return filePathsTimeModified[i].timeModified.Before(filePathsTimeModified[j].timeModified)
	})
	startDateRange := 0
	for i := 0; i < len(filePathsTimeModified)-1; i++ {
		// TODO: duplicate OR i + 1 or i+1
		iPlusOne := i + 1
		if isWithinThreeDays(filePathsTimeModified[iPlusOne].timeModified, filePathsTimeModified[i].timeModified) {
			name, err := createDirectoryDateRangeName(filePathsTimeModified[startDateRange].timeModified, filePathsTimeModified[i].timeModified)
			if err != nil {
				return err
			}
			subDirectoryPath := filepath.Join(destinationDirectory, name)
			foundIndex := -1
			// TODO: add the missing filePathsTimeModified to the dir and remove subDirectoryPath from dateRangeDirectoryPaths in this loop?
			for j, path := range dateRangeDirectoryPaths {
				if path == subDirectoryPath {
					foundIndex = j
					break
				}
			}
			if foundIndex >= 0 {
				dateRangeDirectoryPaths[foundIndex] = dateRangeDirectoryPaths[len(dateRangeDirectoryPaths)-1]
				dateRangeDirectoryPaths = dateRangeDirectoryPaths[:len(dateRangeDirectoryPaths)-1]
			} else {
				if err := utils.CreateDirectory(subDirectoryPath); err != nil {
					return err
				}
			}
			var lastPathElements []string
			for j := startDateRange; j <= i; j++ {
				base := filepath.Base(filePathsTimeModified[j].filePath)
				for _, element := range lastPathElements {
					if element == base {
						return fmt.Errorf("wants to move two files with the same name, '%s' in the directory '%s'", base, subDirectoryPath)
					}
				}
				fullFilePath := filepath.Join(subDirectoryPath, base)
				exists, err := utils.FileOrDirectoryExists(fullFilePath)
				if err != nil {
					return err
				}
				if !exists {
					if err := os.Rename(filePathsTimeModified[j].filePath, fullFilePath); err != nil {
						return err
					}
				}

				// Removing this check from the loop by extracting code to a function can result in an 'err != null' check in this loop since that function can return an error.
				// Therefore, this current check results in less code.
				if j < i {
					lastPathElements = append(lastPathElements, base)
				}
			}
			startDateRange = iPlusOne
		}
	}
	for _, path := range dateRangeDirectoryPaths {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

func formatDateAndWriteString(builder *strings.Builder, time time.Time) error {
	if _, err := builder.WriteString(toDateFormat(time)); err != nil {
		return err
	}
	return nil
}

func createDirectoryFileSystemNodeInSlice(path string) []utils.FileSystemNode {
	return []utils.FileSystemNode{{
		Path:        path,
		IsDirectory: true,
	}}
}

func appendFilePathsTimeModified(filePathsTimeModified *[]filePathTimeModified, uniqueFileSystemNodes []utils.FileSystemNode) error {
	return utils.AppendFileDetails(
		func(detail utils.FileDetail) {
			*filePathsTimeModified = append(*filePathsTimeModified, filePathTimeModified{
				filePath:     detail.Path,
				timeModified: detail.ModificationTime,
			})
		}, uniqueFileSystemNodes, utils.FilesWithoutZeroByteFiles)
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
