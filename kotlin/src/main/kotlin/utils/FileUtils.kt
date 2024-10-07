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
  var hash: String
)

enum class FileFilterMode {
  FILES,
  NON_ZERO_BYTE_FILES,
  FILES_AND_DIRECTORIES,
  NON_ZERO_BYTE_FILES_AND_DIRECTORIES,
  DIRECTORIES
}

enum class FileType {
  ALL_FILES,
  PLAIN_TEXT_FILES
}

fun WalkFilterAndHandleFileMetadata(
  absoluteFilePath: String, mode: FileFilterMode, type: FileType, handler: (FileMetadata?) -> Unit) {
  val rootDirectory = File(absoluteFilePath)

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
    if (file.isDirectory && (mode == FileFilterMode.FILES || mode == FileFilterMode.NON_ZERO_BYTE_FILES)) {
      continue
    }

    // file type check
    if (type == FileType.PLAIN_TEXT_FILES) {
      continue // TODO:
    }

    // is zero byte file check
    if (file.isFile && size == 0L &&
      (mode == FileFilterMode.NON_ZERO_BYTE_FILES || mode == FileFilterMode.NON_ZERO_BYTE_FILES_AND_DIRECTORIES)) {
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
