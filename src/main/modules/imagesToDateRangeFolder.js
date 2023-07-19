import fs from 'fs'
import path from 'path'
import {
  combinePathParts,
  getBaseName,
  getDirectoryFileObjects,
  getDistinctDirectoryFileObjects,
  makeDirectoryIfNotExists,
  removeEmptyDirectories
} from './filePaths.js'
import { filePathsType, fileType } from '../../preload/modules/files'
import {
  isResultObjectOk,
  isResultObjectPartiallyOk,
  toResultObjectWithNullResultAndResultStatusErrorSystem,
  toResultObjectWithNullResultAndResultStatusOk,
  toResultObjectWithNullResultAndResultStatusPartiallyOk,
  toResultObjectWithNullResultByResultObject
} from '../../preload/modules/resultStatus'

// TODO: check for good error handling whole app
// TODO: function looks a lot like duplicateFiles
// TODO: rename resultStatus file
// TODO: check if has access to input en output directory
export default async function imagesToDateRangeFolder(filePaths, outputPath) {
  const inputPath = getSelectedFolderPath(filePaths)

  const imageFilePathsTreeRO = await getDirectoryFileObjects(
    inputPath,
    true,
    filePathsType.filesWithoutZeroByteFiles,
    fileType.image
  )
  if (!isResultObjectOk(imageFilePathsTreeRO)) {
    return toResultObjectWithNullResultByResultObject(imageFilePathsTreeRO)
  }

  const directoryFilePathsRO = await getDirectoryFileObjects(
    outputPath,
    false,
    filePathsType.directories
  )
  if (!isResultObjectOk(directoryFilePathsRO)) {
    return toResultObjectWithNullResultByResultObject(directoryFilePathsRO)
  }

  try {
    const groups = getDateRangeGroups([
      ...imageFilePathsTreeRO.result,
      ...getDateSubdirectoryFileObjects(directoryFilePathsRO.result)
    ])
    await groupsToFolders(groups, outputPath)
  } catch (error) {
    return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
  }

  const removeEmptyDirectoriesRO = await removeEmptyDirectories([
    ...getDistinctDirectoryFileObjects(imageFilePathsTreeRO.result),
    ...directoryFilePathsRO.result
  ])
  if (isResultObjectOk(removeEmptyDirectoriesRO)) {
    return toResultObjectWithNullResultAndResultStatusOk()
  } else if (isResultObjectPartiallyOk(removeEmptyDirectoriesRO)) {
    return toResultObjectWithNullResultAndResultStatusPartiallyOk(removeEmptyDirectoriesRO.message)
  } else {
    return toResultObjectWithNullResultAndResultStatusErrorSystem(removeEmptyDirectoriesRO.message)
  }
}

// TODO: remove function
function getSelectedFolderPath(files) {
  const firstFolderPath = path.dirname(files[0])
  const lastFolderPath = path.dirname(files[files.length - 1])

  let prefix = ''

  for (let i = 0; i < firstFolderPath.length; i++) {
    if (firstFolderPath[i] === lastFolderPath[i]) {
      prefix += firstFolderPath[i]
    } else {
      break
    }
  }

  return prefix
}

function getDateSubdirectoryFileObjects(fileObjects) {
  let result = []
  const separator = ' - '

  for (const fileObject of fileObjects) {
    const baseName = getBaseName(fileObject.path)

    if (baseName.includes(separator)) {
      const directoryParts = baseName.split(separator)
      if (isValidDateFormat(directoryParts[0]) && isValidDateFormat(directoryParts[1])) {
        result = addFileObjects(result, fileObject.path)
      }
    } else if (isValidDateFormat(baseName)) {
      result = addFileObjects(result, fileObject.path)
    }
  }

  return result
}

function addFileObjects(result, fileObjectPath) {
  return [...result, ...getFileObjects(fileObjectPath)]
}

function getDateRangeGroups(fileObjects) {
  fileObjects.sort(compare)

  const groups = []
  let group = [fileObjects[0]]
  for (let i = 1; i < fileObjects.length; i++) {
    const fileObject = fileObjects[i - 1]
    const fileObject2 = fileObjects[i]

    if (isWithinThreeDays(fileObject.dateCreated, fileObject2.dateCreated)) {
      group.push(fileObject2)
    } else {
      groups.push(group)
      group = [fileObject2]
    }
  }
  groups.push(group)

  return groups
}

function isWithinThreeDays(date1, date2) {
  const millisecondsDifference = Math.abs(date2.getTime() - date1.getTime())
  const days = millisecondsDifference / (1000 * 3600 * 24)

  return days <= 3
}

async function groupsToFolders(groups, path) {
  for (const group of groups) {
    const oldestDate = formatTime(group[0].dateCreated)
    const newestDate = formatTime(group[group.length - 1].dateCreated)

    let subFolderPath = combinePathParts(path, oldestDate)
    if (oldestDate !== newestDate) {
      subFolderPath = `${subFolderPath} - ${newestDate}`
    }
    // TODO: error handling
    const makeDirectoryIfNotExistsRO = await makeDirectoryIfNotExists(subFolderPath)
    for (const combination of group) {
      const fileName = getBaseName(combination.path)
      const destinationPath = combinePathParts(subFolderPath, fileName)
      fs.rename(combination.path, destinationPath, (err) => {
        if (err) {
          throw err
        }
      })
    }
  }
}

function isValidDateFormat(dateString) {
  const dateParts = dateString.split('-')

  if (dateParts.length !== 3) {
    return false
  }

  const year = parseInt(dateParts[0], 10)
  const month = parseInt(dateParts[1], 10)
  const day = parseInt(dateParts[2], 10)

  return !(
    isNaN(day) ||
    isNaN(month) ||
    isNaN(year) ||
    day < 1 ||
    day > 31 ||
    month < 1 ||
    month > 12 ||
    year < 0
  )
}

// TODO: use filePaths.js?
function getFileObjects(filePath) {
  const fileObjects = []

  try {
    const files = fs.readdirSync(filePath)

    files.forEach((file) => {
      const combinedFilePath = combinePathParts(filePath, file)
      const stats = fs.statSync(combinedFilePath)

      if (stats.isFile()) {
        // TODO: dateCreated and size correct?
        fileObjects.push({ path: combinedFilePath, dateCreated: 0, size: 0 })
      }
    })
  } catch (error) {
    console.error('Error occurred while reading the folder:', error)
  }

  return fileObjects
}

function formatTime(time) {
  const date = new Date(time)
  const day = date.getDate()
  // Months are zero-based
  const month = date.getMonth() + 1
  const year = date.getFullYear()

  const formattedDay = addPrefix0IfLessThan10(day)
  const formattedMonth = addPrefix0IfLessThan10(month)

  return `${year}-${formattedMonth}-${formattedDay}`
}

function addPrefix0IfLessThan10(number) {
  return number < 10 ? `0${number}` : number
}

function compare(a, b) {
  if (a.dateCreated < b.dateCreated) {
    return -1
  }
  if (a.dateCreated > b.dateCreated) {
    return 1
  }
  return 0
}
