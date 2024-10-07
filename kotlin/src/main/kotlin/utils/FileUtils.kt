package org.example.utils

import java.io.File

data class FileSystemFile(val data: String, val fileMetadata: FileMetadata)

data class FileMetadata(
  val name: String,
  val absoluteDirectoryPath: String,
  val absolutePath: String,
  val timeModified: Long,
  val size: Long,
  val isDirectory: Boolean,
  var hash: String,
)

enum class FileFilterMode {
  FILES,
  NON_ZERO_BYTE_FILES,
  FILES_AND_DIRECTORIES,
  NON_ZERO_BYTE_FILES_AND_DIRECTORIES,
  DIRECTORIES,
}

fun WalkFilterAndHandleFileMetadata(absoluteFilePath: String, mode: FileFilterMode, handler: (FileMetadata?)) {
  val rootDirectory = File("path/to/your/directory")

  for (file in rootDirectory.walk()) {
    var size: Long = 0L
    if (file.isFile) {
      size = file.length()
    }

    // is file check
    if (file.isFile && mode == FileFilterMode.DIRECTORIES) {
      continue
    }

    // is directory check
    if (file.isDirectory && (mode == FileFilterMode.FILES) || mode == FileFilterMode.NON_ZERO_BYTE_FILES) {
      continue
    }

    // is zero byte file check
    if (file.isFile && size == 0L && (mode == FileFilterMode.NON_ZERO_BYTE_FILES || mode == FileFilterMode.NON_ZERO_BYTE_FILES_AND_DIRECTORIES)) {
      continue
    }
  }
}

// TODO: should the Golang version return a FileMetadata{} instead of an error?
fun ToFileMetadata(absoluteFilePath: String): FileMetadata? {
  val file = File(absoluteFilePath)

  if (!file.exists()) {
    return null
  }

  return FileMetadata(
    name = file.name,
    absoluteDirectoryPath = "",
    absolutePath = absoluteFilePath,
    timeModified = file.lastModified(),
    size = file.length(),
    isDirectory = file.isDirectory,
    hash = ""
  )
}
