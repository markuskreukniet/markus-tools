import fs from 'fs'
import path from 'path'
import { inputError } from '../../preload/modules/errors'
import { filePathsType, fileType } from '../../preload/modules/files'
import {
  resultStatus,
  toResultObject,
  toResultObjectWithNullResult
} from '../../preload/modules/resultStatus'

// TODO: change fs import to promises
export async function getDirectoryFilePaths(directoryPath, directoryTree, typeFilePaths, typeFile) {
  if (!typeFile) {
    typeFile = fileType.all
  }

  if (typeFilePaths === filePathsType.directories && typeFile !== fileType.all) {
    return toResultObject([], resultStatus.errorSystem, inputError.wrongFunctionArguments)
  }

  const filePaths = []
  const stack = [directoryPath]
  while (stack.length > 0) {
    const currentPath = stack.pop()

    try {
      const files = await fs.promises.readdir(currentPath)

      const stats = await Promise.all(
        files.map((file) => {
          return fs.promises.stat(path.join(currentPath, file))
        })
      )

      for (let i = 0; i < files.length; i++) {
        const filePath = path.join(currentPath, files[i])

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
      const files = await fs.promises.readdir(filePath)
      if (files.length === 0) {
        await fs.promises.rmdir(filePath)
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

export function getDistinctDirectoryPaths(filePaths) {
  const directoryPaths = filePaths.map((filePath) => path.dirname(filePath))
  const sortedDirectoryPaths = directoryPaths.sort()
  return sortedDirectoryPaths.filter(
    (sortedDirectoryPath, index) => sortedDirectoryPath !== sortedDirectoryPaths[index - 1]
  )
}
