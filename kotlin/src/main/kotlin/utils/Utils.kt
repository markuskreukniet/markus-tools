package org.example.utils

import java.security.MessageDigest

fun createFileMetadataByHashGroups(files: Array<FileMetadata>, onlyDuplicates: Boolean) {
  if (files.isEmpty()) {
    return
  }

  data class FilesByFileSize(
    val fileSize: Long,
    val files: MutableList<FileMetadata>
  )

  fun addGroup(groups: MutableList<FilesByFileSize>, file: FileMetadata) {
    groups.add(FilesByFileSize(
      fileSize = file.size,
      files = mutableListOf(file)
    ))
  }

  val result: MutableList<MutableList<FileMetadata>> = mutableListOf()
  val groups = mutableListOf<FilesByFileSize>()
  var sizeIndex = 0

  files.sortBy { it.size }
  addGroup(groups, files.first())

  files.withIndex().drop(1).forEach { (index, file) ->
    if (file.size == groups[sizeIndex].files.first().size) {
      groups[sizeIndex].files.add(files[index])
    } else {
      addGroup(groups, files[index])
      sizeIndex++
    }
  }

  groups.forEach { group ->
    if (group.files.size > 1) {
      val map = mutableMapOf<String, FileMetadata>()
      group.files.forEach { file ->
        file.absolutePath
        //
      }
      //
    } else if (!onlyDuplicates) {
      result.add(group.files)
    }
  }
}

fun createFileHash(filePath: String): Result<String> {
  val file = createExistingFile(filePath).getOrElse { return Result.failure(it) } ?: return Result.success("")

  val bytes = file.readBytes()
  val md = runCatching { MessageDigest.getInstance("SHA-256") }.getOrElse { return Result.failure(it) }
  val hashBytes = md.digest(bytes)
  return Result.success(hashBytes.joinToString("") { "%02x".format(it) })
}
