package org.example

import java.io.File
import java.nio.file.Path
import java.time.LocalDate
import java.time.format.DateTimeFormatter
import java.time.temporal.ChronoUnit
import kotlin.math.abs

fun isValidDateRangeDirectoryName(name: String): Boolean {
  val spacedHyphen = " - "

  val parseDate = fun(rawDate: String) = runCatching {
    LocalDate.parse(rawDate, DateTimeFormatter.ofPattern("yyyy-MM-dd"))
  }

  if (name.contains(spacedHyphen)) {
    val nameParts = name.split(spacedHyphen)
    val firstDate = parseDate(nameParts.first()).getOrElse { return false }
    val secondDate = parseDate(nameParts[1]).getOrElse { return false }
    if (abs(ChronoUnit.DAYS.between(firstDate, secondDate)) in 1..3) { // TODO: can be longer than three days
      return true
    }
  } else {
    parseDate(name).onSuccess { return true }
  }

  return false
}

fun appendSubtreeContents(directory: File, files: MutableList<File>, directories: MutableList<File>) {
  directory.walk().drop(1).forEach { file ->
    if (file.isDirectory) {
      directories.add(file)
    } else {
      files.add(file)
    }
  }
}

fun categorizeFilesAndDirectories(
  destinationDirectory: File
): Result<Pair<MutableList<File>, Pair<MutableList<File>, MutableList<File>>>> = runCatching {
  val files = mutableListOf<File>()
  val goodDirectories = mutableListOf<File>()
  val badDirectories = mutableListOf<File>()

  // TODO: naming
  val test = fun(directories: MutableList<File>) {
    directories.forEach { directory ->
      appendSubtreeContents(directory, files, badDirectories)
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

  test(goodDirectories)
  test(badDirectories)

  Pair(files, Pair(goodDirectories, badDirectories))
}

fun filesToDateRangeDirectory(uniqueAbsolutePaths: Array<Path>, destinationDirectory: Path) {

}
