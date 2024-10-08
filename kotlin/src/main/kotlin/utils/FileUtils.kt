package org.example.utils

import java.io.File
import java.net.URLConnection

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

data class FileSystemNode(
  val absolutePath: String,
  val isDirectory: Boolean
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
  TEXT_FILES
}

fun isTextFile(file: File): Boolean {
  val mimeType = file.inputStream().use { inputStream ->
    URLConnection.guessContentTypeFromStream(inputStream)
  }

  return mimeType?.startsWith("text") == true
}

fun filterAndHandleFileMetadata(
  file: File, mode: FileFilterMode, type: FileType, absoluteFilePath: String, handler: (FileMetadata) -> Unit) {
  val size: Long = if (file.isFile) file.length() else 0L

  // is file check
  if (file.isFile && mode == FileFilterMode.DIRECTORIES) {
    return
  }

  // is directory check
  if (file.isDirectory && (mode == FileFilterMode.FILES || mode == FileFilterMode.NON_ZERO_BYTE_FILES)) {
    return
  }

  // is zero byte file check
  if (file.isFile && size == 0L &&
    (mode == FileFilterMode.NON_ZERO_BYTE_FILES || mode == FileFilterMode.NON_ZERO_BYTE_FILES_AND_DIRECTORIES)) {
    return
  }

  // is text file check
  if (type == FileType.TEXT_FILES && isTextFile(file)) {
    return
  }

  handler(FileMetadata(
    name = file.name,
    absoluteDirectoryPath = "", // TODO:
    absolutePath = absoluteFilePath,
    timeModified = file.lastModified(),
    size = file.length(),
    isDirectory = file.isDirectory,
    hash = ""
  ))
}

fun walkFilterAndHandleFileMetadata(
  absoluteFilePath: String, mode: FileFilterMode, type: FileType, handler: (FileMetadata) -> Unit) {
  val rootFile = File(absoluteFilePath)

  if (!rootFile.exists()) {
    return
  }

  if (rootFile.isFile) {
    filterAndHandleFileMetadata(rootFile, mode, type, absoluteFilePath, handler)
  } else if (rootFile.isDirectory) {
    rootFile.walk().forEach { file ->
      filterAndHandleFileMetadata(file, mode, type, absoluteFilePath, handler)
    }
  }
}
