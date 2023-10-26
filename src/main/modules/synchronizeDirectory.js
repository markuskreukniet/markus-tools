import {
  combineOutputFilePathWithRelativeInputFilePath,
  copyDirectoryTree,
  copyFile,
  filePathExists,
  getFileAndDirectoryFileObjects
} from './filePaths.js'
import { isResultObjectOk } from '../../preload/modules/resultStatus'

export default async function synchronizeDirectory(
  originalDirectoryFilePath,
  destinationDirectoryFilePath
) {
  const directoriesTree = true
  // added getFileAndDirectoryFileObjects for synchronizeDirectory

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

    // TODO: almost the same as originalFileAndDirectoryFileObjectsRO
    const destinationFileAndDirectoryFileObjectsRO = await getFileAndDirectoryFileObjects(
      destinationDirectoryPath,
      directoriesTree
    )
    if (!isResultObjectOk(destinationFileAndDirectoryFileObjectsRO)) {
      return destinationFileAndDirectoryFileObjectsRO
    }

    for (const fileObject of originalFileAndDirectoryFileObjectsRO.result) {
      // example of fileObject.path: C:\Users\shono\Desktop\test\test\New folder

      const outputFilePath = combineOutputFilePathWithRelativeInputFilePath(
        originalDirectoryFilePath,
        fileObject.path,
        destinationDirectoryFilePath
      )

      const filePathExistsRO = filePathExists(outputFilePath)
      if (!isResultObjectOk(filePathExistsRO)) {
        return filePathExistsRO
      }

      if (filePathExistsRO.result) {
        if (!fileObject.isDirectory) {
          // get output fileObject
          // compare modified date time
          // if destination does have the file and original file is newer, replace the file
          // copyFile does replace, fs.copyFile and fs.createWriteStream, both do that, keep this comment, but in filePaths.js
          await copyFile()
        }
      } else {
        if (fileObject.isDirectory) {
          await copyDirectoryTree()
        } else {
          await copyFile()
        }
      }
    }

    // check for files in destination that are removed in original and remove them from destination
  }

  return `${originalDirectoryFilePath} testB ${destinationDirectoryFilePath}`
}
