// We can't use symbols across the Electron IPC (inter-process communication) boundary
export const filePathsType = Object.freeze({
  files: 'files',
  filesWithoutZeroByteFiles: 'filesWithoutZeroByteFiles',
  filesAndDirectories: 'filesAndDirectories',
  filesAndDirectoriesWithoutZeroByteFiles: 'filesAndDirectoriesWithoutZeroByteFiles',
  directories: 'directories'
})

export const fileType = Object.freeze({
  all: 'all',
  image: 'image'
})

export const filePathSelectionType = Object.freeze({
  both: 'both',
  file: 'file',
  directory: 'directory'
})
