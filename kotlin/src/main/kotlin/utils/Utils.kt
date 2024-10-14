package org.example.utils

import java.security.MessageDigest

fun createFileMetadataByHashGroups(
  files: MutableList<FileMetadata>, onlyDuplicates: Boolean
): Result<MutableList<MutableList<FileMetadata>>?> = runCatching {
  if (files.isEmpty()) {
    return@runCatching null
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
      val map = mutableMapOf<String, MutableList<FileMetadata>>()
      group.files.forEach { file ->
        val hash = createFileHash(file.absolutePath).getOrThrow() ?: return@runCatching null
        map.getOrPut(hash) { mutableListOf() }.add(file)
      }
      map.values.forEach { hashedFiles ->
        if (hashedFiles.size > 1 || !onlyDuplicates){
          result.add(hashedFiles)
        }
      }
    } else if (!onlyDuplicates) {
      result.add(group.files)
    }
  }

  result
}

fun createFileHash(filePath: String): Result<String?> = runCatching {
  val file = createExistingFile(filePath).getOrThrow() ?: return@runCatching null
  val bytes = file.readBytes()
  val md = MessageDigest.getInstance("SHA-256")
  val hashBytes = md.digest(bytes)
  hashBytes.joinToString("") { "%02x".format(it) }
}
