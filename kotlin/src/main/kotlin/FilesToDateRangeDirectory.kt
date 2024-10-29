package org.example

import org.example.utils.*
import java.io.File
import java.nio.file.Path
import java.time.LocalDate
import java.time.format.DateTimeFormatter
import java.time.temporal.ChronoUnit
import kotlin.io.path.exists

fun isValidDateRangeDirectoryName(name: String): Boolean {
  val spacedHyphen = " - "

  val parseDate = fun(rawDate: String) = runCatching {
    LocalDate.parse(rawDate, DateTimeFormatter.ofPattern("yyyy-MM-dd"))
  }

  if (name.contains(spacedHyphen)) {
    val nameParts = name.split(spacedHyphen)
    val firstDate = parseDate(nameParts.first()).getOrElse { return false }
    val secondDate = parseDate(nameParts[1]).getOrElse { return false }
    if (ChronoUnit.DAYS.between(firstDate, secondDate) >= 1) {
      return true
    }
  } else {
    parseDate(name).onSuccess { return true }
  }

  return false
}

fun categorizeFilesAndDirectories(
  destinationDirectory: File
): Pair<MutableList<File>, Pair<MutableList<File>, MutableList<File>>> {
  val files = mutableListOf<File>()
  val goodDirectories = mutableListOf<File>()
  val badDirectories = mutableListOf<File>()

  // TODO: duplicate
//  val handler = fun(file: FileInfo) {
//    if (file.isDirectory) {
//      badDirectories.add(file.file)
//    } else {
//      files.add(file.file)
//    }
//  }

  val categorizeSubtreeContents = fun(directories: MutableList<File>) {
    directories.forEach { directory ->
      directory.walk().drop(1).forEach { file ->
        if (!file.isFile && !file.isDirectory) {
          // exception // duplicate code
        }

        if (file.isDirectory) {
          badDirectories.add(file)
        } else {
          files.add(file)
        }
      }
    }
  }

  destinationDirectory.walk().maxDepth(1).forEach { file ->
    if (!file.isFile && !file.isDirectory) {
      // exception // duplicate code
    }

    if (file.isDirectory) {
      if (isValidDateRangeDirectoryName(file.name)) {
        goodDirectories.add(file)
      } else {
        badDirectories.add(file)
      }
    } else {
      files.add(file)
    }
  }

  categorizeSubtreeContents(goodDirectories)
  categorizeSubtreeContents(badDirectories)

  return Pair(files, Pair(goodDirectories, badDirectories))
}

fun filesToDateRangeDirectory(
  uniqueAbsolutePaths: Array<Path>, destinationDirectory: File
): Result<Unit> = runCatching {
  if (!destinationDirectory.exists()) {
    return@runCatching
  }

  val pair = categorizeFilesAndDirectories(destinationDirectory)
  val files = pair.first
  val goodDirectories = pair.second.first
  val badDirectories = pair.second.second

  // TODO: duplicate
  val handler = fun(file: FileInfo) {
    files.add(file.file)
  }

  uniqueAbsolutePaths.forEach { path ->
    if (path.exists()) {
      walkFilterAndHandleFileInfo(path, FileFilterMode.NON_ZERO_BYTE_FILES, FileType.ALL_FILES, handler)
    }
  }

  // There is no need to check if the directory exists before attempting removal.
  badDirectories.asReversed().forEach { directory ->
    directory.delete()
  }
  goodDirectories.forEach { directory ->
    directory.delete()
  }
}
