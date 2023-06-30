import fs from 'fs'

export async function getDirectoryFilePaths(path, directoryTree, typeFilePaths, typeFileType) {
  if (typeFilePaths === filePathsType.directories && typeFileType !== fileType.all) {
    return []
  }

  const filePaths = []
  const stack = [path]
  while (stack.length > 0) {
    const currentPath = stack.pop()

    try {
      const files = await fs.promises.readdir(currentPath)

      const statsPromises = files.map((file) => {
        const filePath = toFilePath(currentPath, file)
        return fs.promises.stat(filePath)
      })

      const stats = await Promise.all(statsPromises)

      for (let i = 0; i < files.length; i++) {
        const filePath = toFilePath(currentPath, files[i])

        const isDirectory = stats[i].isDirectory()
        if (directoryTree && isDirectory) {
          stack.push(filePath)
        }
        if (shouldAddFilePath(typeFilePaths, typeFileType, filePath, isDirectory, stats.size)) {
          filePaths.push(filePath)
        }
      }
    } catch (err) {
      console.error(err)
    }
  }

  return filePaths
}

function shouldAddFilePath(typeFilePaths, typeFileType, filePath, isDirectory, size) {
  let fileTypeCheck = true
  if (typeFileType === fileType.image && !isDirectory) {
    const lowerCaseFilePath = filePath.toLowerCase()
    fileTypeCheck =
      lowerCaseFilePath.endsWith('jpg') ||
      lowerCaseFilePath.endsWith('jpeg') ||
      lowerCaseFilePath.endsWith('png') ||
      lowerCaseFilePath.endsWith('gif') ||
      lowerCaseFilePath.endsWith('webp')
  }

  const directoryCheck = typeFilePaths === filePathsType.directories && !isDirectory ? false : true
  const zeroByteCheck =
    (typeFilePaths === filePathsType.filesWithoutZeroByteFiles ||
      typeFilePaths === filePathsType.filesAndDirectoriesWithoutZeroByteFiles) &&
    size === 0
      ? false
      : true

  return fileTypeCheck && directoryCheck && zeroByteCheck
}

function toFilePath(path, file) {
  return `${path}\\${file}`
}

export default function isNotAZeroByteFile(stats) {
  return stats.size > 0
}

// We can't use symbols across the Electron IPC (inter-process communication) boundary
const filePathsType = Object.freeze({
  files: 'files',
  filesWithoutZeroByteFiles: 'filesWithoutZeroByteFiles',
  filesAndDirectories: 'filesAndDirectories',
  filesAndDirectoriesWithoutZeroByteFiles: 'filesAndDirectoriesWithoutZeroByteFiles',
  directories: 'directories'
})

const fileType = Object.freeze({
  all: 'all',
  image: 'image'
})
