import { constants, promises } from 'fs'
import path from 'path'
import { ErrorTracker, inputError } from '../../preload/modules/errors'
import { filePathsType, filePathType, fileType } from '../../preload/modules/files'
import {
  isResultObjectOk,
  resultStatus,
  toResultObject,
  toResultObjectWithEmptyArrayResultAndResultStatusErrorSystem,
  toResultObjectWithNullResultAndResultStatusErrorSystem,
  toResultObjectWithNullResultAndResultStatusOk,
  toResultObjectWithResultStatusOk
} from '../../preload/modules/resultStatus'

export async function filePathObjectsToFileObjects(filePathObjects, useDirectoriesTreeInput) {
  const errorTracker = new ErrorTracker()

  const inputImageFileObjects = []
  for (const filePathObject of filePathObjects) {
    let inputRO = null

    if (filePathObject.filePathType === filePathType.file) {
      // TODO: should be getImageFileObject? probably not, only image selection should happen in dialog
      inputRO = await getFileObject(filePathObject.value)
    } else {
      inputRO = await getDirectoryImageFileObjectsWithoutZeroByteOnes(
        filePathObject.value,
        useDirectoriesTreeInput
      )
    }

    if (isResultObjectOk(inputRO)) {
      if (Array.isArray(inputRO.result)) {
        inputImageFileObjects.push(...inputRO.result)
      } else {
        inputImageFileObjects.push(inputRO.result)
      }
    } else {
      errorTracker.concatErrorMessageOnNewLineAndIncrementErrorCount(inputRO.message)
    }
  }

  return errorTracker.createResultObject(inputImageFileObjects.length, inputImageFileObjects)
}

// TODO: use isDirectory param
async function getFileObject(filePath, isDirectory) {
  try {
    const stat = await promises.stat(filePath)
    // TODO: dateCreated or dateModified?
    return toResultObjectWithResultStatusOk({
      path: filePath,
      dateCreated: stat.mtime,
      size: stat.size,
      isDirectory: isDirectory || stat.isDirectory()
    })
  } catch (error) {
    return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
  }
}

// A 'Promise.all' solution for getting all the filePath stats of a folder results in an O(n2) solution.
// Also, with that solution, we can't use the getFileObject function efficiently since that function also gets filePath stats.
// TODO: remove export, also others maybe
export async function getDirectoryFileObjects(
  directoryPath,
  directoryTree,
  typeFilePaths,
  typeFile
) {
  if (!typeFile) {
    typeFile = fileType.all
  }

  if (typeFilePaths === filePathsType.directories && typeFile !== fileType.all) {
    return toResultObjectWithEmptyArrayResultAndResultStatusErrorSystem(
      inputError.wrongFunctionArguments
    )
  }

  const fileObjects = []

  const stack = [directoryPath]
  while (stack.length > 0) {
    const currentPath = stack.pop()

    let files = []
    try {
      files = await promises.readdir(currentPath)
    } catch (error) {
      // TODO: partially ok?
      return toResultObjectWithEmptyArrayResultAndResultStatusErrorSystem(error.message)
    }

    for (const file of files) {
      const filePath = combinePathParts(currentPath, file)
      const fileObjectRO = await getFileObject(filePath)
      // TODO: partially ok?
      if (!isResultObjectOk(fileObjectRO)) {
        return toResultObjectWithEmptyArrayResultAndResultStatusErrorSystem(fileObjectRO.message)
      }

      if (directoryTree && fileObjectRO.result.isDirectory) {
        stack.push(filePath)
      }

      if (
        shouldAddFileObject(
          typeFilePaths,
          typeFile,
          filePath,
          fileObjectRO.result.isDirectory,
          fileObjectRO.result.size
        )
      ) {
        fileObjects.push(fileObjectRO.result)
      }
    }
  }

  return toResultObject(fileObjects, resultStatus.ok)
}

export async function getDirectoryDirectoryFileObjects(directoryPath, directoryTree) {
  return getDirectoryFileObjects(directoryPath, directoryTree, filePathsType.directories)
}

export async function getDirectoryFileObjectsWithoutZeroByteOnes(directoryPath, directoryTree) {
  return getDirectoryFileObjects(
    directoryPath,
    directoryTree,
    filePathsType.filesWithoutZeroByteFiles,
    fileType.all
  )
}

export async function getDirectoryImageFileObjectsWithoutZeroByteOnes(
  directoryPath,
  directoryTree
) {
  return getDirectoryFileObjects(
    directoryPath,
    directoryTree,
    filePathsType.filesWithoutZeroByteFiles,
    fileType.image
  )
}

// TODO: maybe function is useless since objects might not be needed
export async function removeEmptyDirectories(fileObjects) {
  const errorTracker = new ErrorTracker()

  // Both awaits are needed, therefore, a 'await Promise.all' solution is useless.
  for (const fileObject of fileObjects) {
    try {
      const files = await promises.readdir(fileObject.path)
      if (files.length === 0) {
        await promises.rmdir(fileObject.path)
      }
    } catch (error) {
      errorTracker.concatErrorMessageOnNewLineAndIncrementErrorCount(error.message)
    }
  }

  return errorTracker.createResultObject(fileObjects.length)
}

// could create reference problems when using these fileObjects after using this function
export function getDistinctDirectoryFileObjects(fileObjects) {
  for (const fileObject of fileObjects) {
    fileObject.path = path.dirname(fileObject.path)
  }
  // TODO: maybe sorting is unnecessary since all files are read in a directory before going to the next directory.
  fileObjects.sort(compare)
  // TODO: dirty index === 0?
  return fileObjects.filter(
    (fileObject, index) => index === 0 || fileObject.path !== fileObjects[index - 1].path
  )
}

function compare(a, b) {
  if (a.path < b.path) {
    return -1
  }
  if (a.path > b.path) {
    return 1
  }
  return 0
}

function shouldAddFileObject(typeFilePaths, typeFile, filePath, isDirectory, size) {
  const fileTypeCheck =
    isDirectory ||
    typeFile === fileType.all ||
    (typeFile === fileType.image && isImageFilePath(filePath))

  const zeroByteCheck =
    isDirectory ||
    typeFilePaths === filePathsType.files ||
    typeFilePaths === filePathsType.directories ||
    typeFilePaths === filePathsType.filesAndDirectories ||
    ((typeFilePaths === filePathsType.filesWithoutZeroByteFiles ||
      typeFilePaths === filePathsType.filesAndDirectoriesWithoutZeroByteFiles) &&
      size > 0)

  return fileTypeCheck && calculateDirectoryCheck(typeFilePaths, isDirectory) && zeroByteCheck
}

function calculateDirectoryCheck(typeFilePaths, isDirectory) {
  const isDirectoryCheck =
    (typeFilePaths === filePathsType.directories ||
      typeFilePaths === filePathsType.filesAndDirectories ||
      typeFilePaths === filePathsType.filesAndDirectoriesWithoutZeroByteFiles) &&
    isDirectory

  const isNotDirectoryCheck =
    (typeFilePaths === filePathsType.files ||
      typeFilePaths === filePathsType.filesWithoutZeroByteFiles) &&
    !isDirectory

  return isDirectoryCheck || isNotDirectoryCheck
}

function isImageFilePath(filePath) {
  const lowerCaseFilePath = filePath.toLowerCase()
  return (
    lowerCaseFilePath.endsWith('jpg') ||
    lowerCaseFilePath.endsWith('jpeg') ||
    lowerCaseFilePath.endsWith('png') ||
    lowerCaseFilePath.endsWith('gif') ||
    lowerCaseFilePath.endsWith('webp')
  )
}

export function getBaseName(filePath) {
  return path.basename(filePath)
}

export function combinePathParts(filePath1, filePath2) {
  return path.join(filePath1, filePath2)
}

async function filePathExists(filePath) {
  try {
    await promises.access(filePath, constants.F_OK)
    return toResultObjectWithResultStatusOk(true)
  } catch {
    return toResultObjectWithResultStatusOk(false)
  }
}

export async function makeDirectoryIfNotExists(filePath) {
  if (await filePathExists(filePath)) {
    try {
      await promises.mkdir(filePath)
    } catch (error) {
      return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
    }
  }
  return toResultObjectWithNullResultAndResultStatusOk()
}

export async function moveFile(sourcePath, destinationPath) {
  try {
    await promises.rename(sourcePath, destinationPath)
    return toResultObjectWithNullResultAndResultStatusOk()
  } catch (error) {
    return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
  }
}
