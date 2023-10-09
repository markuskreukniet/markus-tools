import { copyDirectoryTree, copyFile, getFileAndDirectoryFileObjects } from './filePaths.js'

export default async function synchronizeDirectory(
  originalDirectoryFilePathObject,
  destinationDirectoryFilePathObject
) {
  // TODO: this boolean should come from UI
  const directoriesTree = true

  // added getFileAndDirectoryFileObjects for synchronizeDirectory
  const originalFileAndDirectoryFileObjectsRO = await getFileAndDirectoryFileObjects(
    originalDirectoryFilePathObject.value,
    directoriesTree
  )
  const destinationFileAndDirectoryFileObjectsRO = await getFileAndDirectoryFileObjects(
    destinationDirectoryFilePathObject.value,
    directoriesTree
  )

  const stack = [originalDirectoryFilePathObject]
  while (stack.length > 0) {
    const currentPath = stack.pop()
  }

  for (const fileObject of originalFileAndDirectoryFileObjectsRO.result) {
    if (fileObject.isDirectory) {
      // if destination does not have the directory, copyDirectoryTree
      await copyDirectoryTree()
      // else stack.push(fileObject.path)
    } else {
      // if destination does not have the file
      await copyFile()
      // if destination does have the file and original file is newer, replace the file
      // TODO: copyFile does replace?
      await copyFile()
    }
  }
}
