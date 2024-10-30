package org.example.utils

import java.io.File
import java.security.MessageDigest

// writeNewline
fun writeNewlineString(builder: StringBuilder) {
  builder.append("\n")
}

fun writeTwoNewlineStrings(builder: StringBuilder) {
  builder.append("\n\n")
}

fun createDuplicateFileInfoGroupsByHash(
  files: MutableList<DuplicateFileInfo>, onlyDuplicates: Boolean
): Result<MutableList<MutableList<DuplicateFileInfo>>?> = runCatching {
  if (files.isEmpty()) {
    return@runCatching null
  }

  data class FilesByFileSize(
    val fileSize: Long,
    val files: MutableList<DuplicateFileInfo>
  )

  val addGroup = fun(groups: MutableList<FilesByFileSize>, file: DuplicateFileInfo) {
    groups.add(FilesByFileSize(
      fileSize = file.size,
      files = mutableListOf(file)
    ))
  }

  val result = mutableListOf<MutableList<DuplicateFileInfo>>()
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
      val map = mutableMapOf<String, MutableList<DuplicateFileInfo>>()
      group.files.forEach { file ->
        val hash = createFileHash(file.file).getOrThrow() ?: return@runCatching null
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

fun createFileHash(file: File): Result<String?> = runCatching {
  val bytes = file.readBytes()
  val md = MessageDigest.getInstance("SHA-256")
  val hashBytes = md.digest(bytes)
  hashBytes.joinToString("") { "%02x".format(it) }
}
