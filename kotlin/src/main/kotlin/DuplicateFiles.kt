package org.example

import org.example.utils.*

fun getDuplicateFilesAsNewlineSeparatedString(uniqueFileSystemNodes: Array<FileSystemNode>): String {
  data class DuplicateFileMetadata(
    val absolutePath: String,
    override val size: Long,
    var hash: String
  ) : FileMetadata

  val result = StringBuilder()
  val files = mutableListOf<DuplicateFileMetadata>()

  val handler = fun(file: CompleteFileMetadata) {
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
