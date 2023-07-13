import fs from 'fs'
import { inputError } from '../../preload/modules/errors'
import { filePathsType, fileType } from '../../preload/modules/files'
import {
  resultStatus,
  toResultObject,
  toResultObjectWithNullResult
} from '../../preload/modules/resultStatus'

export async function getDirectoryFilePaths(path, directoryTree, typeFilePaths, typeFile) {
  if (!typeFile) {
    typeFile = fileType.all
  }

  if (typeFilePaths === filePathsType.directories && typeFile !== fileType.all) {
    return toResultObject([], resultStatus.errorSystem, inputError.wrongFunctionArguments)
  }

  const filePaths = []
  const stack = [path]
  while (stack.length > 0) {
    const currentPath = stack.pop()

    try {
      const files = await fs.promises.readdir(currentPath)

      const statsPromises = files.map((file) => {
        return fs.promises.stat(toFilePath(currentPath, file))
      })
      const stats = await Promise.all(statsPromises)

      for (let i = 0; i < files.length; i++) {
        const filePath = toFilePath(currentPath, files[i])

        const isDirectory = stats[i].isDirectory()
        if (directoryTree && isDirectory) {
          stack.push(filePath)
        }
        if (shouldAddFilePath(typeFilePaths, typeFile, filePath, isDirectory, stats.size)) {
          filePaths.push(filePath)
        }
      }
    } catch (error) {
      return toResultObject([], resultStatus.errorSystem, error.message)
    }
  }

  return toResultObject(filePaths, resultStatus.ok)
}

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

function toFilePath(path, file) {
  return `${path}\\${file}`
}

// TODO: remove export
export default function isNotAZeroByteFile(stats) {
  return stats.size > 0
}

export async function removeEmptyDirectories(filePaths) {
  let errorCount = 0
  let errorMessage = ''

  // Both awaits are needed, therefore, a 'await Promise.all' solution is useless.
  for (const path of filePaths) {
    try {
      const files = await fs.promises.readdir(path)
      if (files.length === 0) {
        await fs.promises.rmdir(path)
      }
    } catch (error) {
      errorCount++
      errorMessage = `${errorMessage}\n${error.message}`
    }
  }

  if (errorCount === 0) {
    return toResultObjectWithNullResult(resultStatus.ok)
  } else if (errorCount > 0 && errorCount < filePaths.length) {
    return toResultObjectWithNullResult(resultStatus.partiallyOk, errorMessage)
  } else {
    return toResultObjectWithNullResult(resultStatus.errorSystem, errorMessage)
  }
}

// TODO: change fs import to promises
