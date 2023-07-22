import { constants, promises } from 'fs'
import path from 'path'
import { ErrorTracker, inputError } from '../../preload/modules/errors'
import { filePathsType, fileType } from '../../preload/modules/files'
import {
  resultStatus,
  toResultObject,
  toResultObjectWithEmptyArrayResultAndResultStatusErrorSystem,
  toResultObjectWithNullResultAndResultStatusErrorSystem,
  toResultObjectWithNullResultAndResultStatusOk,
  toResultObjectWithNullResultAndResultStatusPartiallyOk,
  toResultObjectWithResultStatusOk
} from '../../preload/modules/resultStatus'

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

    try {
      const files = await promises.readdir(currentPath)

      const stats = await Promise.all(
        files.map((file) => {
          return promises.stat(combinePathParts(currentPath, file))
        })
      )

      for (let i = 0; i < files.length; i++) {
        const filePath = combinePathParts(currentPath, files[i])

        const isDirectory = stats[i].isDirectory()
        if (directoryTree && isDirectory) {
          stack.push(filePath)
        }
        if (shouldAddFileObject(typeFilePaths, typeFile, filePath, isDirectory, stats[i].size)) {
          // TODO: dateCreated or dateModified?
          fileObjects.push({
            path: filePath,
            dateCreated: stats[i].mtime,
            size: stats[i].size
          })
        }
      }
    } catch (error) {
      return toResultObjectWithEmptyArrayResultAndResultStatusErrorSystem(error.message)
    }
  }

  return toResultObject(fileObjects, resultStatus.ok)
}

export async function getDirectoryDirectoryFileObjects(directoryPath, directoryTree) {
  return getDirectoryFileObjects(directoryPath, directoryTree, filePathsType.directories)
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
  // TODO: use it
  const errorTracker = new ErrorTracker()

  // TODO: error object with update function?
  let errorCount = 0
  let errorMessage = ''

  // Both awaits are needed, therefore, a 'await Promise.all' solution is useless.
  for (const fileObject of fileObjects) {
    try {
      const files = await promises.readdir(fileObject.path)
      if (files.length === 0) {
        await promises.rmdir(fileObject.path)
      }
    } catch (error) {
      errorCount++
      errorMessage = `${errorMessage}\n${error.message}`
    }
  }

  if (errorCount === 0) {
    return toResultObjectWithNullResultAndResultStatusOk()
  } else if (errorCount > 0 && errorCount < fileObjects.length) {
    return toResultObjectWithNullResultAndResultStatusPartiallyOk(errorMessage)
  } else {
    return toResultObjectWithNullResultAndResultStatusErrorSystem(errorMessage)
  }
}

export function getDistinctDirectoryFileObjects(fileObjects) {
  for (const fileObject of fileObjects) {
    fileObject.path = path.dirname(fileObject.path)
  }
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

// TODO: remove export
export default function isNotAZeroByteFile(stats) {
  return stats.size > 0
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
      return toResultObjectWithNullResultAndResultStatusOk()
    } catch (error) {
      return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
    }
  } else {
    return false
  }
}

export async function moveFile(sourcePath, destinationPath) {
  try {
    await promises.rename(sourcePath, destinationPath)
    return toResultObjectWithNullResultAndResultStatusOk()
  } catch (error) {
    return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
  }
}
