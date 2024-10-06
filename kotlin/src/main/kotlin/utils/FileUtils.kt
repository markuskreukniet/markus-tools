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
    timeModified = file.lastModified(), // TODO: format time
    size = file.length(),
    isDirectory = file.isDirectory(),
    hash = ""
  )
}
