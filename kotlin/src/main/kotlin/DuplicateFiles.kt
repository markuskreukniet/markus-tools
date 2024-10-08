package org.example

import org.example.utils.*

data class DuplicateFileMetadata(
  val absolutePath: String,
  val size: Long,
  var hash: String
)

fun getDuplicateFilesAsNewlineSeparatedString(uniqueFileSystemNodes: Array<FileSystemNode>): String {
  val result = StringBuilder()
  val files = mutableListOf<DuplicateFileMetadata>()

  val handler = fun(file: FileMetadata) {
    files.add(DuplicateFileMetadata(
      absolutePath = file.absolutePath,
      size = file.size,
      hash = file.hash
    ))
  }

  uniqueFileSystemNodes.forEach { node ->
    walkFilterAndHandleFileMetadata(node.absolutePath, FileFilterMode.NON_ZERO_BYTE_FILES, FileType.ALL_FILES, handler)
  }

  result.append("test")

  return result.toString()
}
