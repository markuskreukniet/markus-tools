import { copyDirectoryTree, copyFile, getFileAndDirectoryFileObjects } from './filePaths.js'

export default async function synchronizeDirectory(
  originalDirectoryFilePathObject,
  destinationDirectoryFilePathObject
) {
  return `${originalDirectoryFilePathObject.value} testB ${destinationDirectoryFilePathObject.value}`

  // TODO: this boolean should come from UI
  const directoriesTree = true
  // added getFileAndDirectoryFileObjects for synchronizeDirectory

  const stack = [originalDirectoryFilePathObject.value]
  while (stack.length > 0) {
    const originalDirectoryPath = stack.pop()
    const destinationDirectoryPath = originalDirectoryPath.replace(
      originalDirectoryFilePathObject.value,
      destinationDirectoryFilePathObject.value
    )

    const originalFileAndDirectoryFileObjectsRO = await getFileAndDirectoryFileObjects(
      originalDirectoryPath,
      directoriesTree
    )
    const destinationFileAndDirectoryFileObjectsRO = await getFileAndDirectoryFileObjects(
      destinationDirectoryPath,
      directoriesTree
    )

    for (const fileObject of originalFileAndDirectoryFileObjectsRO.result) {
      if (fileObject.isDirectory) {
        // if destination does not have the directory, copyDirectoryTree
        await copyDirectoryTree()
        // else stack.push(fileObject.path)
      } else {
        // if destination does not have the file
        await copyFile()
        // if destination does have the file and original file is newer, replace the file
        // copyFile does replace, fs.copyFile and fs.createWriteStream, both do that, keep this comment, but in filePaths.js
        await copyFile()
      }
    }

    // check for files in destination that are removed in original and remove them from destination
  }
}
