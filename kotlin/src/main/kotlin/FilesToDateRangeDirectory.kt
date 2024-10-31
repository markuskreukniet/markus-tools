package org.example

import org.example.utils.DuplicateFileInfo
import org.example.utils.FTDRFileInfo
import org.example.utils.createDuplicateFileInfoGroupsByHash
import java.io.File
import java.nio.file.Path
import java.time.LocalDate
import java.time.format.DateTimeFormatter
import java.time.temporal.ChronoUnit
import kotlin.io.path.exists
import kotlin.io.path.getLastModifiedTime

val addDirectory = fun(directories: MutableList<File>, file: File) {
  directories.add(file)
}

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

  val categorizeInDirectory = fun(directories: MutableList<File>, file: File) {
    if (isValidDateRangeDirectoryName(file.name)) {
      goodDirectories.add(file)
    } else {
      directories.add(file)
    }
  }

  val categorizeSubtreeContents = fun(directories: MutableList<File>) {
    directories.forEach { directory ->
      directory.walk().drop(1).forEach { file ->
        categorize(file, files, badDirectories, addDirectory)
      }
    }
  }

  destinationDirectory.walk().maxDepth(1).forEach { file ->
    categorize(file, files, badDirectories, categorizeInDirectory)
  }

  categorizeSubtreeContents(goodDirectories)
  categorizeSubtreeContents(badDirectories)

  return Pair(files, Pair(goodDirectories, badDirectories))
}

fun categorize(
  file: File,
  files: MutableList<File>,
  badDirectories: MutableList<File>,
  handler: (directories: MutableList<File>, file: File) -> Unit
) {
  if (file.isDirectory) {
    handler(badDirectories, file)
  } else if (file.isFile) {
    if (file.length() > 0L) {
      files.add(file)
    } else {
      // exception
    }
  } else {
    // exception
  }
}

fun deleteDuplicateFiles(files: MutableList<FTDRFileInfo>, destinationDirectory: File) = runCatching {
  val groups = createDuplicateFileInfoGroupsByHash(files, false).getOrThrow() ?: return@runCatching null

}

fun filesToDateRangeDirectory(
  uniqueAbsolutePaths: Array<Path>, destinationDirectory: File
): Result<Unit> = runCatching {
  if (!destinationDirectory.exists()) {
    return@runCatching
  }

  // TODO: duplicate values
  val pair = categorizeFilesAndDirectories(destinationDirectory)
  val files = pair.first
  val goodDirectories = pair.second.first
  val badDirectories = pair.second.second

  uniqueAbsolutePaths.forEach { path ->
    if (path.exists()) {
      categorize(path.toFile(), files, badDirectories, addDirectory)
    }
  }

  // TODO: remove this converting
  val files2 = mutableListOf<FTDRFileInfo>()
  files.forEach { file ->
    val absolutePath = file.toPath().toAbsolutePath()
    files2.add(FTDRFileInfo(
      file = file,
      size = file.length(),
      absolutePath = absolutePath,
      timeModified = absolutePath.getLastModifiedTime()
    ))
  }

  // delete duplicate files

  files2.sortBy { it.timeModified }

  // There is no need to check if the directory exists before attempting removal.
  badDirectories.asReversed().forEach { directory ->
    directory.delete()
  }
  goodDirectories.forEach { directory ->
    directory.delete()
  }
}
