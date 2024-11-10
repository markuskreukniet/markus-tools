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

fun <T : DuplicateFileInfo> createDuplicateFileInfoGroupsByHash(
  files: MutableList<T>, onlyDuplicates: Boolean
): Result<MutableList<MutableList<T>>?> = runCatching {
  if (files.isEmpty()) {
    return@runCatching null
  }

  data class FilesByFileSize(
    val fileSize: Long,
    val files: MutableList<T>
  )

  val result = mutableListOf<MutableList<T>>()
  val groups = mutableListOf<FilesByFileSize>()
  var groupIndex = 0

  val addGroup = fun(file: T) {
    groups.add(FilesByFileSize(
      fileSize = file.size,
      files = mutableListOf(file)
    ))
  }

  files.sortBy { it.size }
  addGroup(files.first())

  files.withIndex().drop(1).forEach { (index, file) ->
    if (file.size == groups[groupIndex].files.first().size) {
      groups[groupIndex].files.add(files[index])
    } else {
      addGroup(files[index])
      groupIndex++
    }
  }

  groups.forEach { group ->
    if (group.files.size > 1) {
      val map = mutableMapOf<String, MutableList<T>>()
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
