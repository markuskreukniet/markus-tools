import {
  combineOutputFilePathWithRelativeInputFilePath,
  copyDirectoryTree,
  copyFile,
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

    for (const [
      originalFilePath,
      originalFileObject
    ] of originalFileAndDirectoryFileObjectsPathMap.entries()) {
      // example of fileObject.path: C:\Users\shono\Desktop\test\test\New folder

      const outputFilePath = combineOutputFilePathWithRelativeInputFilePath(
        originalDirectoryFilePath,
        originalFilePath,
        destinationDirectoryFilePath
      )

      if (destinationFileAndDirectoryFileObjectsPathMap.has(outputFilePath)) {
        if (originalFileObject.isDirectory) {
          stack.push(originalFilePath)
        } else {
          if (
            originalFileObject.dateModified >
            destinationFileAndDirectoryFileObjectsPathMap.get(outputFilePath).dateModified
          ) {
            // TODO: RO
            // copyFile does replace, fs.copyFile and fs.createWriteStream, both do that, keep this comment, but in filePaths.js
            await copyFile(originalFilePath, outputFilePath)
          }
        }
      } else {
        if (originalFileObject.isDirectory) {
          // TODO: RO
          await copyDirectoryTree()
        } else {
          // TODO: RO and copied
          await copyFile(originalFilePath, outputFilePath)
        }
      }

      for (const [
        destinationFilePath,
        destinationFileObject
      ] of destinationFileAndDirectoryFileObjectsPathMap.entries()) {
        if (!originalFileAndDirectoryFileObjectsPathMap.has(destinationFilePath)) {
          if (destinationFileObject.isDirectory) {
            // TODO: RO
            await removeDirectoryTree(destinationFilePath)
          } else {
            // TODO: RO
            await removeFile(destinationFilePath)
          }
        }
      }
    }
  }

  return `${originalDirectoryFilePath} testB ${destinationDirectoryFilePath}`
}
