package org.example

import org.example.utils.*
import java.nio.file.Path

fun plainTextFilesToText(uniqueAbsolutePaths: Array<Path>): String? {
  val files = mutableListOf<FTextFilesFileInfo>()

  val handler = fun(file: CompleteFileInfo) {
    files.add(
      FTextFilesFileInfo(
        file = file.file,
        absolutePath = file.absolutePath
      )
    )
  }

  uniqueAbsolutePaths.forEach { path ->
    walkFilterAndHandleFileInfo(path, FileFilterMode.NON_ZERO_BYTE_FILES, FileType.TEXT_FILES, handler)
  }

  if (files.isEmpty()) {
    return null
  }

  val result = StringBuilder()

  return  result.toString()
}
