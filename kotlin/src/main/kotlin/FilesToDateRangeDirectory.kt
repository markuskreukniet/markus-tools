package org.example

import org.example.utils.FTDRFileInfo
import org.example.utils.createDuplicateFileInfoGroupsByHash
import java.io.File
import java.nio.file.Path
import java.time.LocalDate
import java.time.ZoneId
import java.time.format.DateTimeFormatter
import java.time.temporal.ChronoUnit
import kotlin.io.path.exists
import kotlin.io.path.getLastModifiedTime

fun getDateTimeFormatter(): Result<DateTimeFormatter> = runCatching {
  DateTimeFormatter.ofPattern("yyyy-MM-dd")
}

val addDirectory = fun(directories: MutableList<File>, file: File) {
  directories.add(file)
}

fun isValidDateRangeDirectoryName(name: String): Boolean {
  val spacedHyphen = " - "

  val parseDate = fun(rawDate: String) = runCatching {
    LocalDate.parse(rawDate, getDateTimeFormatter().getOrThrow())
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
): Pair<MutableList<File>, Pair<MutableMap<String, File>, MutableList<File>>> {
  val files = mutableListOf<File>()
  val goodDirectoriesByName = mutableMapOf<String, File>() // TODO: is the File needed?
  val badDirectories = mutableListOf<File>()

  val categorizeInDirectory = fun(directories: MutableList<File>, file: File) {
    if (isValidDateRangeDirectoryName(file.name)) {
      goodDirectoriesByName[file.name] = file
    } else {
      directories.add(file)
    }
  }

  val categorizeSubtreeContents = fun(directories: MutableCollection<File>) {
    directories.forEach { directory ->
      directory.walk().drop(1).forEach { file ->
        categorize(file, files, badDirectories, addDirectory)
      }
    }
  }

  destinationDirectory.walk().maxDepth(1).forEach { file ->
    categorize(file, files, badDirectories, categorizeInDirectory)
  }

  categorizeSubtreeContents(goodDirectoriesByName.values)
  categorizeSubtreeContents(badDirectories)

  return Pair(files, Pair(goodDirectoriesByName, badDirectories))
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

fun createHandlers(
  destinationDirectory: File
): List<(MutableList<FTDRFileInfo>, MutableList<File>) -> MutableList<FTDRFileInfo>> {
  val addAllFilesInfo = fun(files: MutableList<File>, filesInfo: List<FTDRFileInfo>) {
    filesInfo.forEach { file ->
      files.add(file.file)
    }
  }

  // TODO: naming
  val addBadFilesInfoAndReplaceGoodFile  = fun(
    badFiles: MutableList<File>, goodFiles: MutableList<FTDRFileInfo>, file: FTDRFileInfo
  ) {
    addAllFilesInfo(badFiles, goodFiles)
    goodFiles.clear()
    goodFiles.add(file)
  }

  val categorizeOnShortestFileNameLength = fun(
    files: MutableList<FTDRFileInfo>, badFiles: MutableList<File>
  ): MutableList<FTDRFileInfo> {
    val good = mutableListOf(files.first())
    var minimumLength = files.first().file.name.length

    files.drop(1).forEach { file ->
      if (file.file.name.length < minimumLength) {
        minimumLength = file.file.name.length
        addBadFilesInfoAndReplaceGoodFile(badFiles, good, file)
      } else if (file.file.name.length == minimumLength) {
        good.add(file)
      } else {
        badFiles.add(file.file)
      }
    }

    return good
  }

  val categorizeOnValidDateRangeDirectoryName = fun(
    files: MutableList<FTDRFileInfo>, badFiles: MutableList<File>
  ): MutableList<FTDRFileInfo> {
    val tempGood1Files = mutableListOf<FTDRFileInfo>()
    val tempGood2Files = mutableListOf<FTDRFileInfo>()
    val tempBadFiles = mutableListOf<FTDRFileInfo>()

    files.forEach { file ->
      if (file.file.parentFile.parentFile == destinationDirectory) {
        if (isValidDateRangeDirectoryName(file.file.parentFile.name)) {
          tempGood2Files.add(file)
        } else {
          tempGood1Files.add(file)
        }
      } else {
        tempBadFiles.add(file)
      }
    }

    if (tempGood2Files.isNotEmpty()) {
      addAllFilesInfo(badFiles, tempGood1Files)
      addAllFilesInfo(badFiles, tempBadFiles)
      return tempGood2Files
    }

    if (tempGood1Files.isNotEmpty()) {
      addAllFilesInfo(badFiles, tempBadFiles)
      return tempGood1Files
    }

    return tempBadFiles
  }

  val categorizeOnNewestTimeModified = fun(
    files: MutableList<FTDRFileInfo>, badFiles: MutableList<File>
  ): MutableList<FTDRFileInfo> {
    val good = mutableListOf(files.first())
    var newest = files.first().timeModified

    files.drop(1).forEach { file ->
      if (file.timeModified > newest) {
        newest = file.timeModified
        addBadFilesInfoAndReplaceGoodFile(badFiles, good, file)
      } else if (file.timeModified == newest) {
        good.add(file)
      } else {
        badFiles.add(file.file)
      }
    }

    return good
  }

  val categorizeOnFirstFile = fun(
    files: MutableList<FTDRFileInfo>, badFiles: MutableList<File>
  ): MutableList<FTDRFileInfo> {
    val good = mutableListOf(files.first())

    addAllFilesInfo(badFiles, files.drop(1))

    return good
  }

  return listOf(
    categorizeOnShortestFileNameLength,
    categorizeOnValidDateRangeDirectoryName,
    categorizeOnNewestTimeModified,
    categorizeOnFirstFile
  )
}

fun deleteDuplicateFiles(
  files: MutableList<FTDRFileInfo>, destinationDirectory: File
): Result<MutableList<FTDRFileInfo>?> = runCatching {
  val groups = createDuplicateFileInfoGroupsByHash(files, false).getOrThrow() ?: return@runCatching null
  val handlers = createHandlers(destinationDirectory)
  val badFiles = mutableListOf<File>()

  files.clear()

  groups.forEachIndexed { index, group ->
    for (handler in handlers) {
      // group and groups[index] are different references
      if (groups[index].size > 1) {
        groups[index] = handler(group, badFiles)
      } else {
        files.add(group.first())
        break
      }
    }
  }

  badFiles.forEach { file ->
    file.delete()
  }

  files
}

fun moveFilesToDirectories(files: MutableList<FTDRFileInfo>, goodDirectoriesByName: MutableMap<String, File>) { //runCatching?
  if (files.size == 0) {
    return
  }

  // val toFormattedString(iets: iets): String = runCatching {
  //   iets.atZone(ZoneId.systemDefault())
  //   .toLocalDate()
  //   .format(getDateTimeFormatter().getOrThrow())
  // }

  files.sortBy { it.timeModified }

  var group = mutableListOf(files.first())

  files.drop(1).forEach { file ->
    if (ChronoUnit.DAYS.between(group.last().timeModified.toInstant(), file.timeModified.toInstant()) in 0..3) {
      group.add(file)
    } else {
      // val result = if (condition) valueIfTrue else valueIfFalse

      // var ding: String
      // if (ChronoUnit.DAYS.between(group.last().timeModified.toInstant(), file.timeModified.toInstant()) >= 1) {

      // }

      // group.first()
      // group.last()

      // val ding = file.timeModified.toInstant()
      //   .atZone(ZoneId.systemDefault())
      //   .toLocalDate()
      //   .format(getDateTimeFormatter().getOrThrow())

      // Remove if the key exists
      goodDirectoriesByName.remove(file.file.parentFile.name)

      // move files to dir

      group = mutableListOf(file)
    }
  }

  if (group.size > 0) {

  }
}

fun filesToDateRangeDirectory(
  uniqueAbsolutePaths: Array<Path>, destinationDirectory: File
): Result<Unit> = runCatching {
  if (!destinationDirectory.exists()) {
    return@runCatching
  }

  val pair = categorizeFilesAndDirectories(destinationDirectory)
  val files = pair.first
  val goodDirectoriesByName = pair.second.first
  val badDirectories = pair.second.second

  uniqueAbsolutePaths.forEach { path ->
    if (path.exists()) {
      categorize(path.toFile(), files, badDirectories, addDirectory)
    }
  }

  // TODO: remove this converting
  var files2 = mutableListOf<FTDRFileInfo>()
  files.forEach { file ->
    val absolutePath = file.toPath().toAbsolutePath()
    files2.add(FTDRFileInfo(
      file = file,
      size = file.length(),
      absolutePath = absolutePath,
      timeModified = absolutePath.getLastModifiedTime() // TODO: should be toInstant() or something similar (so not first to getLastModifiedTime)?
    ))
  }

  files2 = deleteDuplicateFiles(files2, destinationDirectory).getOrThrow() ?: return@runCatching

  //

  // There is no need to check if the directory exists before attempting removal.
  badDirectories.asReversed().forEach { directory ->
    directory.delete()
  }
  goodDirectoriesByName.values.forEach { directory ->
    directory.delete()
  }
}
