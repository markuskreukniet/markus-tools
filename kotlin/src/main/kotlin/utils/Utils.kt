package org.example.utils

fun createFileMetadataByHashGroups(files: Array<FileMetadata>, onlyDuplicates: Boolean) {
  if (files.isEmpty()) {
    return
  }

  data class FilesByFileSize(
    val fileSize: Long,
    val files: MutableList<FileMetadata>
  )

  files.sortBy { it.size }

  val result: MutableList<MutableList<FileMetadata>> = mutableListOf()
  val groups = mutableListOf<FilesByFileSize>()
  var sizeIndex = 0

  // TODO: abstract, now duplicate with line 29
  groups.add(FilesByFileSize(
    fileSize = files.first().size,
    files = mutableListOf(files.first())
  ))

  files.withIndex().drop(1).forEach { (index, file) ->
    if (file.size == groups[sizeIndex].files.first().size) {
      groups[sizeIndex].files.add(files[index])
    } else {
      groups.add(FilesByFileSize(
        fileSize = files[index].size,
        files = mutableListOf(files[index])
      ))
      sizeIndex++
    }
  }

  groups.forEach { group ->
    if (group.files.size > 1) {
      //
    } else if (!onlyDuplicates) {
      result.add(group.files)
    }
  }
}
