import { ErrorTracker } from '../../preload/modules/errors'
import {
  combinePathParts,
  filePathObjectsToFileObjects,
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
  toResultObjectWithNullResultAndResultStatusPartiallyOk
} from '../../preload/modules/resultStatus'

// TODO: rename resultStatus file
// TODO: rename useDirectoriesTreeInput, without use?
// TODO: remove param useDirectoriesTreeInput
// TODO: write access check?
export default async function imagesToDateRangeFolder(
  filePathObjects,
  outputPath,
  useDirectoriesTreeInput
) {
  const filePathObjectsToFileObjectsRO = await filePathObjectsToFileObjects(
    filePathObjects,
    useDirectoriesTreeInput
  )
  if (!isResultObjectOk(filePathObjectsToFileObjectsRO)) {
    return filePathObjectsToFileObjectsRO
  }

  const outputImageFileObjectsRO = await getDirectoryImageFileObjectsWithoutZeroByteOnes(
    outputPath,
    false
  )
  if (!isResultObjectOk(outputImageFileObjectsRO)) {
    return outputImageFileObjectsRO
  }

  // TODO: use inputDirectoryFileObjectsTreeRO
  // const inputDirectoryFileObjectsTreeRO = await getDirectoryDirectoryFileObjects(inputPath, true)
  // if (!isResultObjectOk(inputDirectoryFileObjectsTreeRO)) {
  //   return inputDirectoryFileObjectsTreeRO
  // }

  const outputDirectoryFileObjectsRO = await getDirectoryDirectoryFileObjects(outputPath, false)
  if (!isResultObjectOk(outputDirectoryFileObjectsRO)) {
    return outputDirectoryFileObjectsRO
  }

  const groups = getDateRangeGroups([
    ...filePathObjectsToFileObjectsRO.result,
    ...getDateSubdirectoryFileObjects(outputImageFileObjectsRO.result)
  ])

  const groupsToDirectoriesRO = await groupsToDirectories(groups, outputPath)
  if (!isResultObjectOk(groupsToDirectoriesRO)) {
    return groupsToDirectoriesRO
  }

  // the array can have duplicate directories
  // TODO: duplicates might be the problem
  const removeEmptyDirectoriesRO = await removeEmptyDirectories([
    ...getDistinctDirectoryFileObjects(filePathObjectsToFileObjectsRO.result),
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
  // This function should not return [[undefined]] when fileObjects is [].
  if (fileObjects.length === 0) {
    return []
  }

  fileObjects.sort(compare)

  const groups = []
  let group = [fileObjects[0]]
  for (let i = 1; i < fileObjects.length; i++) {
    const fileObject = fileObjects[i - 1]
    const fileObject2 = fileObjects[i]

    if (isWithinThreeDays(fileObject.dateModified, fileObject2.dateModified)) {
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

async function groupsToDirectories(groups, outputPath) {
  const errorTracker = new ErrorTracker()

  for (const group of groups) {
    errorTracker.addNumberOfPossibleErrors(group.length)

    const oldestDate = formatTime(group[0].dateModified)
    const newestDate = formatTime(group[group.length - 1].dateModified)

    let subFolderPath = combinePathParts(outputPath, oldestDate)
    if (oldestDate !== newestDate) {
      subFolderPath = `${subFolderPath} - ${newestDate}`
    }

    const makeDirectoryIfNotExistsRO = await makeDirectoryIfNotExists(subFolderPath)
    if (isResultObjectOk(makeDirectoryIfNotExistsRO)) {
      await Promise.all(
        group.map(async (fileObject) => {
          const moveFileRO = await moveFile(
            fileObject.path,
            combinePathParts(subFolderPath, getBaseName(fileObject.path))
          )

          if (!isResultObjectOk(moveFileRO)) {
            errorTracker.concatErrorMessageOnNewLineAndIncrementErrorCount(moveFileRO.message)
          }
        })
      )
    } else {
      errorTracker.concatErrorMessageOnNewLineAndIncrementErrorCount(
        makeDirectoryIfNotExistsRO.message
      )
    }
  }

  return errorTracker.createResultObject()
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
  if (a.dateModified < b.dateModified) {
    return -1
  }
  if (a.dateModified > b.dateModified) {
    return 1
  }
  return 0
}
