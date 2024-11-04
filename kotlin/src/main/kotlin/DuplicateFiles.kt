package org.example

import org.example.utils.*
import java.nio.file.Path
import kotlin.io.path.exists

fun getDuplicateFilesAsNewlineSeparatedString(
  uniqueAbsolutePaths: Array<Path>
): Result<String?> = runCatching {
  val result = StringBuilder()
  val files = mutableListOf<DuplicateFileInfo>()

  val handler = fun(file: FileInfo) {
    files.add(DFFileInfo(
      file = file.file,
      size = file.size,
      absolutePath = file.absolutePath // TODO: why not use file.file.absolutePath? // TODO: also update DFFileInfo then?
    ))
  }

  uniqueAbsolutePaths.forEach { path ->
    if (path.exists()) {
      walkFilterAndHandleFileInfo(path, FileFilterMode.NON_ZERO_BYTE_FILES, FileType.ALL_FILES, handler)
    }
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
      result.append(file.absolutePath)
    }
  }

  result.toString()
}
