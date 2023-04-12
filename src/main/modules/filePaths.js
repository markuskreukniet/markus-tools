// async function getFileObject(filePath, isDirectory) {
//   try {
//     const stat = await promises.stat(filePath)
//     return toResultObjectWithResultStatusOk({
//       path: filePath,
//       dateModified: stat.mtime,
//       size: stat.size,
//       isDirectory: isDirectory || stat.isDirectory()
//     })
//   } catch (error) {
//     return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
//   }
// }

// export async function removeDirectoryTree(filePath) {
//   try {
//     await promises.rm(filePath, { recursive: true })
//     return toResultObjectWithNullResultAndResultStatusOk()
//   } catch (error) {
//     return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
//   }
// }

// export async function removeFile(filePath) {
//   try {
//     await promises.rm(filePath)
//     return toResultObjectWithNullResultAndResultStatusOk()
//   } catch (error) {
//     return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
//   }
// }

// await promises.rmdir(fileObject.path)

// export function getBaseName(filePath) {
//   return path.basename(filePath)
// }

// export function combinePathParts(filePath1, filePath2) {
//   return path.join(filePath1, filePath2)
// }

// export async function filePathExists(filePath) {
//   try {
//     await promises.access(filePath, constants.F_OK)
//     return toResultObjectWithResultStatusOk(true)
//   } catch {
//     return toResultObjectWithResultStatusOk(false)
//   }
// }

// async function readFilesFromDirectory(filePath) {
//   try {
//     return toResultObjectWithResultStatusOk(await promises.readdir(filePath))
//   } catch (error) {
//     return toResultObjectWithEmptyArrayResultAndResultStatusErrorSystem(error.message)
//   }
// }

// export async function makeDirectory(filePath) {
//   try {
//     await promises.mkdir(filePath)
//     return toResultObjectWithNullResultAndResultStatusOk()
//   } catch (error) {
//     return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
//   }
// }

// export async function moveFile(sourcePath, destinationPath) {
//   try {
//     await promises.rename(sourcePath, destinationPath)
//     return toResultObjectWithNullResultAndResultStatusOk()
//   } catch (error) {
//     return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
//   }
// }

// TODO: promises.copyFile does not work efficient with huge files
// export async function copyFile(sourcePath, destinationPath) {
//   try {
//     await promises.copyFile(sourcePath, destinationPath)
//     return toResultObjectWithNullResultAndResultStatusOk()
//   } catch (error) {
//     return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
//   }
// }

// export async function getReadFileHandle(filePath) {
//   try {
//     return toResultObjectWithResultStatusOk(await open(filePath, 'r'))
//   } catch (error) {
//     return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
//   }
// }

// function getRelativePath(filePathFrom, filePathTo) {
//   return path.relative(filePathFrom, filePathTo)
// }
