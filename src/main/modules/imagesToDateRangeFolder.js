import path from 'path'
import { ErrorTracker } from '../../preload/modules/errors'
import {
  combinePathParts,
  getBaseName,
  getDirectoryDirectoryFileObjects,
  getDirectoryImageFileObjectsWithoutZeroByteOnes,
  getDistinctDirectoryFileObjects,
  makeDirectoryIfNotExists,
  moveFile,
  removeEmptyDirectories
} from './filePaths.js'
import {
  isResultObjectOk,
  isResultObjectPartiallyOk,
  toResultObjectWithNullResultAndResultStatusErrorSystem,
  toResultObjectWithNullResultAndResultStatusOk,
  toResultObjectWithNullResultAndResultStatusPartiallyOk,
  toResultObjectWithNullResultByResultObject
} from '../../preload/modules/resultStatus'

// TODO: check for good error handling whole app
// TODO: rename resultStatus file
// TODO: check if has access to input en output directory
// TODO: remove path import
export default async function imagesToDateRangeFolder(filePaths, outputPath) {
  const inputPath = getSelectedFolderPath(filePaths)

  const inputImageFileObjectsTreeRO = await getDirectoryImageFileObjectsWithoutZeroByteOnes(
    inputPath,
    true
  )
  if (!isResultObjectOk(inputImageFileObjectsTreeRO)) {
    return toResultObjectWithNullResultByResultObject(inputImageFileObjectsTreeRO)
  }

  const outputImageFileObjectsRO = await getDirectoryImageFileObjectsWithoutZeroByteOnes(
    outputPath,
    false
  )
  if (!isResultObjectOk(outputImageFileObjectsRO)) {
    return toResultObjectWithNullResultByResultObject(outputImageFileObjectsRO)
  }

  // TODO: use inputDirectoryFileObjectsTreeRO
  const inputDirectoryFileObjectsTreeRO = await getDirectoryDirectoryFileObjects(inputPath, true)
  if (!isResultObjectOk(inputDirectoryFileObjectsTreeRO)) {
    return toResultObjectWithNullResultByResultObject(inputDirectoryFileObjectsTreeRO)
  }

  const outputDirectoryFileObjectsRO = await getDirectoryDirectoryFileObjects(outputPath, false)
  if (!isResultObjectOk(outputDirectoryFileObjectsRO)) {
    return toResultObjectWithNullResultByResultObject(outputDirectoryFileObjectsRO)
  }

  const groups = getDateRangeGroups([
    ...inputImageFileObjectsTreeRO.result,
    ...getDateSubdirectoryFileObjects(outputImageFileObjectsRO.result)
  ])

  // TODO: error handling
  await groupsToDirectories(groups, outputPath)

  // the array can have duplicate directories
  // TODO: duplicates might be the problem
  const removeEmptyDirectoriesRO = await removeEmptyDirectories([
    ...getDistinctDirectoryFileObjects(inputImageFileObjectsTreeRO.result),
    // ...getDistinctDirectoryFileObjects(outputImageFileObjectsRO.result),
    // ...inputDirectoryFileObjectsTreeRO.result,
    ...outputDirectoryFileObjectsRO.result
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
  const result = []
  const separator = ' - '

  for (const fileObject of fileObjects) {
    const baseName = getBaseName(fileObject.path)

    if (baseName.includes(separator)) {
      const directoryParts = baseName.split(separator)
      if (isValidDateFormat(directoryParts[0]) && isValidDateFormat(directoryParts[1])) {
        result.push(fileObject)
      }
    } else if (isValidDateFormat(baseName)) {
      result.push(fileObject)
    }
  }

  return result
}

function getDateRangeGroups(fileObjects) {
  // TODO: it might be possible to remove a sort since getDistinctDirectoryFileObjects has also a sort
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

// TODO: the only function left to check/fix, also naming in this function
async function groupsToDirectories(groups, path) {
  const errorTracker = new ErrorTracker()
  let maxPossibleErrors = 0

  for (const group of groups) {
    maxPossibleErrors = maxPossibleErrors + group.length

    const oldestDate = formatTime(group[0].dateCreated)
    const newestDate = formatTime(group[group.length - 1].dateCreated)

    let subFolderPath = combinePathParts(path, oldestDate)
    if (oldestDate !== newestDate) {
      subFolderPath = `${subFolderPath} - ${newestDate}`
    }

    const makeDirectoryIfNotExistsRO = await makeDirectoryIfNotExists(subFolderPath)
    if (isResultObjectOk(makeDirectoryIfNotExistsRO)) {
      for (const combination of group) {
        const moveFileRO = moveFile(
          combination.path,
          combinePathParts(subFolderPath, getBaseName(combination.path))
        )
        if (!isResultObjectOk(moveFileRO)) {
          errorTracker.concatErrorMessageOnNewLineAndIncrementErrorCount(moveFileRO.message)
        }
      }
    } else {
      errorTracker.concatErrorMessageOnNewLineAndIncrementErrorCount(
        makeDirectoryIfNotExistsRO.message
      )
    }
  }

  return errorTracker.toResultObjectWithNullResult(maxPossibleErrors)
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
