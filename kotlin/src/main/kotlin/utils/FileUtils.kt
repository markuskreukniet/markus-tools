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
  var mimeType: String? = null

  file.inputStream().use { inputStream ->
    mimeType = URLConnection.guessContentTypeFromStream(inputStream)
  }

  return mimeType?.startsWith("text") == true
}

fun walkFilterAndHandleFileMetadata(
  absoluteFilePath: String, mode: FileFilterMode, type: FileType, handler: (FileMetadata?) -> Unit) {
  val rootDirectory = File(absoluteFilePath)

  for (file in rootDirectory.walk()) {
    val size: Long = if (file.isFile) file.length() else 0L

    // is file check
    if (file.isFile && mode == FileFilterMode.DIRECTORIES) {
      continue
    }

    // is directory check
    if (file.isDirectory && (mode == FileFilterMode.FILES || mode == FileFilterMode.NON_ZERO_BYTE_FILES)) {
      continue
    }

    // is zero byte file check
    if (file.isFile && size == 0L &&
      (mode == FileFilterMode.NON_ZERO_BYTE_FILES || mode == FileFilterMode.NON_ZERO_BYTE_FILES_AND_DIRECTORIES)) {
      continue
    }

    // is text file check
    if (type == FileType.TEXT_FILES && isTextFile(file)) {
      continue
    }
  }
}

// TODO: should the Golang version return a FileMetadata{} instead of an error?
fun toFileMetadata(absoluteFilePath: String): FileMetadata? {
  val file = File(absoluteFilePath)

  if (!file.exists()) {
    return null
  }

  return FileMetadata(
    name = file.name,
    absoluteDirectoryPath = "", // TODO:
    absolutePath = absoluteFilePath,
    timeModified = file.lastModified(),
    size = file.length(),
    isDirectory = file.isDirectory,
    hash = ""
  )
}
