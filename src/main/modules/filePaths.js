import { constants, promises } from 'fs'
import { open } from 'fs/promises'
import path from 'path'
import { ErrorTracker, inputError } from '../../preload/modules/errors'
import { filePathsType, filePathType, fileType } from '../../preload/modules/files'
import {
  isResultObjectOk,
  toResultObjectWithEmptyArrayResultAndResultStatusErrorSystem,
  toResultObjectWithNullResultAndResultStatusErrorSystem,
  toResultObjectWithNullResultAndResultStatusOk,
  toResultObjectWithResultStatusOk
} from '../../preload/modules/resultStatus'

export async function filePathObjectsToFileObjects(filePathObjects, useDirectoriesTreeInput) {
  const errorTracker = new ErrorTracker(filePathObjects.length)

  // TODO: always ImageFileObjects?
  const inputImageFileObjects = []
  for (const filePathObject of filePathObjects) {
    let inputRO = null

    if (filePathObject.filePathType === filePathType.file) {
      // TODO: should be getImageFileObject? probably not, only image selection should happen in dialog
      inputRO = await getFileObject(filePathObject.value, false)
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

  return errorTracker.createResultObject(inputImageFileObjects)
}

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
      inputError.aWrongCombinationOfArguments
    )
  }

  const errorTracker = new ErrorTracker()
  const fileObjects = []

  const stack = [directoryPath]
  while (stack.length > 0) {
    errorTracker.addNumberOfPossibleErrors(1)

    const currentPath = stack.pop()
    const readFilesFromDirectoryRO = await readFilesFromDirectory(currentPath)
    if (isResultObjectOk(readFilesFromDirectoryRO)) {
      errorTracker.addNumberOfPossibleErrors(readFilesFromDirectoryRO.result.length)

      for (const file of readFilesFromDirectoryRO.result) {
        const filePath = combinePathParts(currentPath, file)
        const fileObjectRO = await getFileObject(filePath)
        if (isResultObjectOk(fileObjectRO)) {
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
        } else {
          errorTracker.concatErrorMessageOnNewLineAndIncrementErrorCount(fileObjectRO.message)
        }
      }
    } else {
      errorTracker.concatErrorMessageOnNewLineAndIncrementErrorCount(
        readFilesFromDirectoryRO.message
      )
    }
  }

  return errorTracker.createResultObject(fileObjects)
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

export async function getFileAndDirectoryFileObjects(directoryPath, directoryTree) {
  return getDirectoryFileObjects(
    directoryPath,
    directoryTree,
    filePathsType.filesAndDirectories,
    fileType.all
  )
}

export async function removeDirectoryTree(filePath) {
  try {
    await promises.rm(filePath, { recursive: true })
    return toResultObjectWithNullResultAndResultStatusOk()
  } catch (error) {
    return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
  }
}

export async function removeFile(filePath) {
  try {
    await promises.rm(filePath)
    return toResultObjectWithNullResultAndResultStatusOk()
  } catch (error) {
    return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
  }
}

// TODO: maybe function is useless since objects might not be needed
// TODO: should use removeDirectoryTree with RO
export async function removeEmptyDirectories(fileObjects) {
  const errorTracker = new ErrorTracker(fileObjects.length)

  // Both awaits are needed, therefore, a 'await Promise.all' solution is useless.
  for (const fileObject of fileObjects) {
    const readFilesFromDirectoryRO = await readFilesFromDirectory(fileObject.path)
    if (isResultObjectOk(readFilesFromDirectoryRO)) {
      try {
        if (readFilesFromDirectoryRO.result.length === 0) {
          await promises.rmdir(fileObject.path)
        }
      } catch (error) {
        errorTracker.concatErrorMessageOnNewLineAndIncrementErrorCount(error.message)
      }
    } else {
      errorTracker.concatErrorMessageOnNewLineAndIncrementErrorCount(
        readFilesFromDirectoryRO.message
      )
    }
  }

  return errorTracker.createResultObject()
}

// TODO: check if true: could create reference problems when using these fileObjects after using this function
export function getDistinctDirectoryFileObjects(fileObjects) {
  for (const fileObject of fileObjects) {
    fileObject.path = path.dirname(fileObject.path)
  }
  // TODO: maybe sorting is unnecessary since all files are read in a directory before going to the next directory.
  fileObjects.sort(compare)

  let previousPath = ''
  return fileObjects.filter((fileObject) => {
    const isUnique = fileObject.path !== previousPath
    previousPath = fileObject.path
    return isUnique
  })
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
      typeFilePaths === filePathsType.filesAndDirectories ||
      typeFilePaths === filePathsType.filesWithoutZeroByteFiles) &&
    !isDirectory

  return isDirectoryCheck || isNotDirectoryCheck
}

// TODO: if you can select only a image, this becomes useless
function isImageFilePath(filePath) {
  const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp']
  return imageExtensions.includes(filePath.toLowerCase().slice(filePath.lastIndexOf('.')))
}

export function getBaseName(filePath) {
  return path.basename(filePath)
}

export function combinePathParts(filePath1, filePath2) {
  return path.join(filePath1, filePath2)
}

export async function filePathExists(filePath) {
  try {
    await promises.access(filePath, constants.F_OK)
    return toResultObjectWithResultStatusOk(true)
  } catch {
    return toResultObjectWithResultStatusOk(false)
  }
}

async function readFilesFromDirectory(filePath) {
  try {
    return toResultObjectWithResultStatusOk(await promises.readdir(filePath))
  } catch (error) {
    return toResultObjectWithEmptyArrayResultAndResultStatusErrorSystem(error.message)
  }
}

export async function makeDirectoryIfNotExists(filePath) {
  if (await filePathExists(filePath)) {
    const makeDirectoryRO = await makeDirectory(filePath)
    if (!isResultObjectOk(makeDirectoryRO)) {
      return makeDirectoryRO
    }
  }
  return toResultObjectWithNullResultAndResultStatusOk()
}

export async function makeDirectory(filePath) {
  try {
    await promises.mkdir(filePath)
    return toResultObjectWithNullResultAndResultStatusOk()
  } catch (error) {
    return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
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

// TODO: promises.copyFile does not work efficient with huge files
export async function copyFile(sourcePath, destinationPath) {
  try {
    await promises.copyFile(sourcePath, destinationPath)
    return toResultObjectWithNullResultAndResultStatusOk()
  } catch (error) {
    return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
  }
}

export async function getReadFileHandle(filePath) {
  try {
    return toResultObjectWithResultStatusOk(await open(filePath, 'r'))
  } catch (error) {
    return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
  }
}

export async function getUtf8FileContents(filePath) {
  try {
    return toResultObjectWithResultStatusOk(await promises.readFile(filePath, { encoding: 'utf8' }))
  } catch (error) {
    return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
  }
}

function getRelativePath(filePathFrom, filePathTo) {
  return path.relative(filePathFrom, filePathTo)
}

// inputFilePathFull.replace to determine the result ends up in a bug.
export function combineOutputFilePathWithRelativeInputFilePath(
  inputFilePathPart,
  inputFilePathFull,
  outputFilePath
) {
  return combinePathParts(outputFilePath, getRelativePath(inputFilePathPart, inputFilePathFull))
}

// TODO useless now?
// When one error happens, the directory tree does not get copied and should return an error.
// TODO: check if copyDirectoryTree is a correct name
export async function copyDirectoryTree(inputFilePath, outputFilePath) {
  // TODO: makeDirectory can fail
  // TODO: check if await makeDirectory(outputFilePath) is needed
  // await makeDirectory(outputFilePath)

  const stack = [inputFilePath]
  while (stack.length > 0) {
    const currentPath = stack.pop()
    const readFilesFromDirectoryRO = await readFilesFromDirectory(currentPath)
    if (!isResultObjectOk(readFilesFromDirectoryRO)) {
      return readFilesFromDirectoryRO
    }

    // TODO: this for can be optimized with promise.all?, maybe also other for loops then
    for (const file of readFilesFromDirectoryRO.result) {
      const inputFilePathFull = combinePathParts(currentPath, file)
      const outputFilePathFull = combineOutputFilePathWithRelativeInputFilePath(
        inputFilePath,
        inputFilePathFull,
        outputFilePath
      )
      const fileObjectRO = await getFileObject(inputFilePathFull)
      if (!isResultObjectOk(fileObjectRO)) {
        return fileObjectRO
      }

      if (fileObjectRO.result.isDirectory) {
        stack.push(inputFilePathFull)
        const makeDirectoryRO = await makeDirectory(outputFilePathFull)
        if (!isResultObjectOk(makeDirectoryRO)) {
          return makeDirectoryRO
        }
      } else {
        await copyFile(inputFilePathFull, outputFilePathFull)
      }
    }
  }

  return toResultObjectWithNullResultAndResultStatusOk()
}
