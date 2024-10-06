package org.example.utils

data class FileSystemFile(val data: String, val fileMetadata: FileMetadata)

data class FileMetadata(
  val name: String,
  val directoryPath: String,
  val path: String,
  val timeModified: Long,
  val size: Long,
  val isDirectory: Boolean,
  var hash: String,
)
