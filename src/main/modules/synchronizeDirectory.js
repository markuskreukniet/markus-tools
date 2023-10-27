import {
  combineOutputFilePathWithRelativeInputFilePath,
  copyDirectoryTree,
  copyFile,
  filePathExists,
  getFileAndDirectoryFileObjects,
  removeDirectoryTree,
  removeFile
} from './filePaths.js'
import { isResultObjectOk } from '../../preload/modules/resultStatus'

export default async function synchronizeDirectory(
  originalDirectoryFilePath,
  destinationDirectoryFilePath
) {
  const directoriesTree = true

  const stack = [originalDirectoryFilePath]
  while (stack.length > 0) {
    const originalDirectoryPath = stack.pop()
    const destinationDirectoryPath = originalDirectoryPath.replace(
      originalDirectoryFilePath,
      destinationDirectoryFilePath
    )

    const originalFileAndDirectoryFileObjectsRO = await getFileAndDirectoryFileObjects(
      originalDirectoryPath,
      directoriesTree
    )
    if (!isResultObjectOk(originalFileAndDirectoryFileObjectsRO)) {
      return originalFileAndDirectoryFileObjectsRO
    }

    // TODO: originalFileAndDirectoryFileObjectsRO.result should be already a Map
    const originalFileAndDirectoryFileObjectsPathMap = new Map(
      originalFileAndDirectoryFileObjectsRO.result.map((fileObject) => [
        fileObject.path,
        {
          dateModified: fileObject.dateCreated,
          isDirectory: fileObject.isDirectory
        }
      ])
    )

    // TODO: almost the same as originalFileAndDirectoryFileObjectsRO
    const destinationFileAndDirectoryFileObjectsRO = await getFileAndDirectoryFileObjects(
      destinationDirectoryPath,
      directoriesTree
    )
    if (!isResultObjectOk(destinationFileAndDirectoryFileObjectsRO)) {
      return destinationFileAndDirectoryFileObjectsRO
    }

    // TODO: originalFileAndDirectoryFileObjectsRO.result should be already a Map
    const destinationFileAndDirectoryFileObjectsPathMap = new Map(
      destinationFileAndDirectoryFileObjectsRO.result.map((fileObject) => [
        fileObject.path,
        {
          dateModified: fileObject.dateCreated,
          isDirectory: fileObject.isDirectory
        }
      ])
    )

    for (const fileObject of originalFileAndDirectoryFileObjectsRO.result) {
      // example of fileObject.path: C:\Users\shono\Desktop\test\test\New folder

      const outputFilePath = combineOutputFilePathWithRelativeInputFilePath(
        originalDirectoryFilePath,
        fileObject.path,
        destinationDirectoryFilePath
      )

      // TODO: not needed anymore? Also remove export in filePaths
      const filePathExistsRO = filePathExists(outputFilePath)
      if (!isResultObjectOk(filePathExistsRO)) {
        return filePathExistsRO
      }

      if (filePathExistsRO.result) {
        if (fileObject.isDirectory) {
          stack.push(fileObject.path)
        } else {
          if (
            fileObject.dateCreated >
            destinationFileAndDirectoryFileObjectsPathMap.get(outputFilePath).dateModified
          ) {
            // TODO: RO
            // copyFile does replace, fs.copyFile and fs.createWriteStream, both do that, keep this comment, but in filePaths.js
            await copyFile(fileObject.path, outputFilePath)
          }
        }
      } else {
        if (fileObject.isDirectory) {
          // TODO: RO
          await copyDirectoryTree()
        } else {
          // TODO: RO and copied
          await copyFile(fileObject.path, outputFilePath)
        }
      }
    }

    if (
      originalFileAndDirectoryFileObjectsPathMap.size !==
      destinationFileAndDirectoryFileObjectsPathMap.size
    ) {
      for (const destinationFileObject of destinationFileAndDirectoryFileObjectsRO.result) {
        // destinationFileObject.path to original and use it in the if
        if (!originalFileAndDirectoryFileObjectsPathMap.has()) {
          if (destinationFileObject.isDirectory) {
            // TODO: RO
            await removeDirectoryTree(destinationFileObject.path)
          } else {
            // TODO: RO
            await removeFile(destinationFileObject.path)
          }
        }
      }
    }
  }

  return `${originalDirectoryFilePath} testB ${destinationDirectoryFilePath}`
}
