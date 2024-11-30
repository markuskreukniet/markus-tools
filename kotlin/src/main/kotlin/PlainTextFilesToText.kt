package org.example

import org.example.utils.*
import java.nio.file.Path

fun plainTextFilesToText(uniqueAbsolutePaths: Array<Path>): Result<String?> = runCatching {
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
    return@runCatching null
  }

  val result = StringBuilder()

  val addNameAndLines = fun(file: FTextFilesFileInfo) {
    result.append(file.file.name)

    file.file.forEachLine { line ->
      writeNewlineString(result)
      result.append(line)
    }
  }

  addNameAndLines(files.first())

  files.drop(1).forEach { file ->
    writeTwoNewlineStrings(result)
    addNameAndLines(file)
  }

  result.toString()
}
