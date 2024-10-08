package org.example.utils

fun createFileMetadataByHashGroups(files: Array<FileMetadata>) {
  if (files.isEmpty()) {
    return
  }

  data class FilesByFileSize(
    val fileSize: Long,
    val files: MutableList<FileMetadata>
  )

  files.sortBy { it.size }
}
