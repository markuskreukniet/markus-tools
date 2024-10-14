package org.example

import org.example.utils.*

fun getDuplicateFilesAsNewlineSeparatedString(
  uniqueFileSystemNodes: Array<FileSystemNode>
): Result<String?> = runCatching {
  data class DuplicateFileMetadata(
    override val absolutePath: String,
    override val size: Long,
    var hash: String
  ) : FileMetadata

  val result = StringBuilder()
  val files = mutableListOf<FileMetadata>()

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

  val groups = createFileMetadataByHashGroups(files, true).getOrThrow() ?: return@runCatching null

  groups.forEachIndexed { indexI, group ->
    if (indexI > 0) {

    }
    group.forEachIndexed { indexJ, file ->
      if (indexJ > 0) {

      }
      //result.append(file.absolutePath)
    }
  }

  result.append("test")

  result.toString()
}
