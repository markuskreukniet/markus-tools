function getDirectoryTreeFilePaths(path, typeFilePaths, typeFileType) {
  if (typeFilePaths === filePathsType.directories && typeFileType !== fileType.all) {
    return []
  }
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
