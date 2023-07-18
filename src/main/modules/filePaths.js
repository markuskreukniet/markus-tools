import { constants, promises } from 'fs'
import path from 'path'
import { inputError } from '../../preload/modules/errors'
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

// new
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
        if (shouldAddFilePath(typeFilePaths, typeFile, filePath, isDirectory, stats[i].size)) {
          // TODO: dateCreated or dateModified?
          fileObjects.push({ path: filePath, dateCreated: stats[i].mtime, size: stats[i].size })
        }
      }
    } catch (error) {
      return toResultObjectWithEmptyArrayResultAndResultStatusErrorSystem(error.message)
    }
  }

  return toResultObject(fileObjects, resultStatus.ok)
}

// rename removeEmptyDirectoriesNew to removeEmptyDirectories
export async function removeEmptyDirectoriesNew(fileObjects) {
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

// old
export async function getDirectoryFilePaths(directoryPath, directoryTree, typeFilePaths, typeFile) {
  if (!typeFile) {
    typeFile = fileType.all
  }

  if (typeFilePaths === filePathsType.directories && typeFile !== fileType.all) {
    return toResultObjectWithEmptyArrayResultAndResultStatusErrorSystem(
      inputError.wrongFunctionArguments
    )
  }

  const filePaths = []
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
        if (shouldAddFilePath(typeFilePaths, typeFile, filePath, isDirectory, stats[i].size)) {
          filePaths.push(filePath)
        }
      }
    } catch (error) {
      return toResultObjectWithEmptyArrayResultAndResultStatusErrorSystem(error.message)
    }
  }

  return toResultObject(filePaths, resultStatus.ok)
}

// TODO: rename shouldAddObject
function shouldAddFilePath(typeFilePaths, typeFile, filePath, isDirectory, size) {
  const fileTypeCheck =
    isDirectory ||
    typeFile === fileType.all ||
    (typeFile === fileType.image && isImageFilePath(filePath))

  const zeroByteCheck =
    (typeFilePaths === filePathsType.filesWithoutZeroByteFiles ||
      typeFilePaths === filePathsType.filesAndDirectoriesWithoutZeroByteFiles) &&
    size === 0
      ? false
      : true

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

export async function removeEmptyDirectories(filePaths) {
  let errorCount = 0
  let errorMessage = ''

  // Both awaits are needed, therefore, a 'await Promise.all' solution is useless.
  for (const filePath of filePaths) {
    try {
      const files = await promises.readdir(filePath)
      if (files.length === 0) {
        await promises.rmdir(filePath)
      }
    } catch (error) {
      errorCount++
      errorMessage = `${errorMessage}\n${error.message}`
    }
  }

  if (errorCount === 0) {
    return toResultObjectWithNullResultAndResultStatusOk()
  } else if (errorCount > 0 && errorCount < filePaths.length) {
    return toResultObjectWithNullResultAndResultStatusPartiallyOk(errorMessage)
  } else {
    return toResultObjectWithNullResultAndResultStatusErrorSystem(errorMessage)
  }
}

export function getBaseName(filePath) {
  return path.basename(filePath)
}

export function combinePathParts(filePath1, filePath2) {
  return path.join(filePath1, filePath2)
}

export function getDistinctDirectoryPaths(filePaths) {
  const directoryPaths = filePaths.map((filePath) => path.dirname(filePath))
  const sortedDirectoryPaths = directoryPaths.sort()
  return sortedDirectoryPaths.filter(
    (sortedDirectoryPath, index) => sortedDirectoryPath !== sortedDirectoryPaths[index - 1]
  )
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
