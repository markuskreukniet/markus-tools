package org.example

import java.io.File
import java.nio.file.Path
import java.time.LocalDate
import java.time.format.DateTimeFormatter
import java.time.temporal.ChronoUnit

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

  val categorizeSubtreeContents = fun(directories: MutableList<File>) {
    directories.forEach { directory ->
      directory.walk().drop(1).forEach { file ->
        if (file.isDirectory) {
          badDirectories.add(file)
        } else {
          files.add(file)
        }
      }
    }
  }

  destinationDirectory.walk().maxDepth(1).forEach { file ->
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

fun filesToDateRangeDirectory(uniqueAbsolutePaths: Array<Path>, destinationDirectory: File) {
  val pair = categorizeFilesAndDirectories(destinationDirectory)
  val files = pair.first
  val goodDirectories = pair.second.first
  val badDirectories = pair.second.second
}
