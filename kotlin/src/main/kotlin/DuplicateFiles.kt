package org.example

import org.example.utils.*

fun getDuplicateFilesAsNewlineSeparatedString(uniqueFileSystemNodes: Array<FileSystemNode>): String {
  data class DuplicateFileMetadata(
    override val absolutePath: String,
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

  // TODO: createFileMetadataByHashGroups moet runCatching worden en mutableList ontangen ipv array
  // val groups = createFileMetadataByHashGroups(files, true).getOrThrow()

  // groups.forEachIndexed { indexI, group ->
  //   if indexI > 0 {

  //   }
  //   group.forEachIndexed { indexJ, file ->
  //     if indexJ > 0 {

  //     }
  //   }
  // }

  result.append("test")

  return result.toString()
}
