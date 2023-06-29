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

        if (stats[i].isDirectory()) {
          stack.push(filePath)
        } else {
          filePaths.push(filePath)
        }
      }
    } catch (err) {
      console.error(err)
    }
  }

  return filePaths
}

function toFilePath(path, file) {
  return `${path}\\${file}`
}

export default function isNotAZeroByteFile(stats) {
  return stats.size > 0
}

// We can't use symbols across the Electron IPC (inter-process communication) boundary
const filePathsType = Object.freeze({
  files: 'ok',
  filesWithoutZeroByteFiles: 'errorSystem',
  filesAndDirectoriesWithoutZeroByteFiles: '',
  directories: ''
})

const fileType = Object.freeze({
  all: 'ok',
  image: 'errorSystem'
})
