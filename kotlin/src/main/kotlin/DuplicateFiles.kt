package org.example

import org.example.utils.*
import java.nio.file.Path

fun getDuplicateFilesAsNewlineSeparatedString(
  uniqueAbsolutePaths: Array<Path>
): Result<String?> = runCatching {
  val result = StringBuilder()
  val files = mutableListOf<DuplicateFileInfo>()

  val handler = fun(file: CompleteFileInfo) {
    files.add(
      FDuplicateFilesFileInfo(
        file = file.file,
        size = file.size
      )
    )
  }

  uniqueAbsolutePaths.forEach { path ->
    walkFilterAndHandleFileInfo(path, FileFilterMode.NON_ZERO_BYTE_FILES, FileType.ALL_FILES, handler)
  }

  val groups = createDuplicateFileInfoGroupsByHash(files, true).getOrThrow() ?: return@runCatching null

  groups.forEachIndexed { indexI, group ->
    if (indexI > 0) {
      writeTwoNewlineStrings(result)
    }
    group.forEachIndexed { indexJ, file ->
      if (indexJ > 0) {
        writeNewlineString(result)
      }
      result.append(file.file.absolutePath)
    }
  }

  result.toString()
}
