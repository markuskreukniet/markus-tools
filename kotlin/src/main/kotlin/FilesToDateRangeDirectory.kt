package org.example

import org.example.utils.FTDRFileInfo
import org.example.utils.createDuplicateFileInfoGroupsByHash
import java.io.File
import java.nio.file.Files
import java.nio.file.Path
import java.nio.file.Paths
import java.time.Instant
import java.time.LocalDate
import java.time.ZoneId
import java.time.format.DateTimeFormatter
import java.time.temporal.ChronoUnit
import kotlin.io.path.exists
import kotlin.io.path.getLastModifiedTime

// TODO: check runCatching for every function

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
): Pair<MutableList<FTDRFileInfo>, Pair<MutableSet<File>, MutableList<File>>> {
  val files = mutableListOf<FTDRFileInfo>()
  val goodDirectories = mutableSetOf<File>()
  val badDirectories = mutableListOf<File>()

  val categorizeInDirectory = fun(directories: MutableList<File>, file: File) {
    if (isValidDateRangeDirectoryName(file.name)) {
      goodDirectories.add(file)
    } else {
      directories.add(file)
    }
  }

  destinationDirectory.walk().maxDepth(1).forEach { file ->
    categorize(file, files, badDirectories, categorizeInDirectory)
  }

  val directories: MutableCollection<File> = mutableListOf()
  directories.addAll(goodDirectories)
  directories.addAll(badDirectories)

  directories.forEach { directory ->
    directory.walk().drop(1).forEach { file ->
      categorize(file, files, badDirectories, addDirectory)
    }
  }

  return Pair(files, Pair(goodDirectories, badDirectories))
}

fun categorize(
  file: File,
  files: MutableList<FTDRFileInfo>,
  badDirectories: MutableList<File>,
  handler: (MutableList<File>, File) -> Unit
) = runCatching {
  if (file.isDirectory) {
    handler(badDirectories, file)
  } else if (file.isFile) {
    val size = file.length()
    if (size > 0L) {
      val absolutePath = file.toPath().toAbsolutePath()
      files.add(FTDRFileInfo(
        file = file,
        size = size,
        absolutePath = absolutePath,
        timeModified = absolutePath.getLastModifiedTime().toInstant()
      ))
    } else {
      // TODO: exception
    }
  } else {
    // TODO: exception
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

  val addBadFilesInfoAndReplaceGoodFiles = fun(
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
        addBadFilesInfoAndReplaceGoodFiles(badFiles, good, file)
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
        addBadFilesInfoAndReplaceGoodFiles(badFiles, good, file)
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
): Result<Unit> = runCatching {
  val groups = createDuplicateFileInfoGroupsByHash(files, false).getOrThrow() ?: return@runCatching
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
}

fun moveFilesAndFilterGoodDirectories(
  files: MutableList<FTDRFileInfo>, goodDirectories: MutableSet<File>, destinationDirectory: File
) = runCatching {
  if (files.isEmpty()) {
    return@runCatching
  }

  files.sortBy { it.timeModified }

  var group = mutableListOf(files.first())

  val toFormattedString = fun(time: Instant): Result<String> = runCatching {
    time.atZone(ZoneId.systemDefault())
      .toLocalDate()
      .format(getDateTimeFormatter().getOrThrow())
  }

  val moveFilesToDirectory = fun() {
    val firstFile = group.first()
    val lastFile = group.last()

    var directoryName = toFormattedString(firstFile.timeModified).getOrThrow()
    if (ChronoUnit.DAYS.between(firstFile.timeModified, lastFile.timeModified) >= 1) {
      directoryName += " - ${toFormattedString(lastFile.timeModified).getOrThrow()}"
    }

    val directory = File(destinationDirectory.absolutePath, directoryName)
    if (directory in goodDirectories) {
      goodDirectories.remove(directory)
    } else {
      directory.mkdir()
    }

    group.forEach { file ->
      val filePath = Paths.get(directory.toString(), file.file.name)
      if (filePath != file.absolutePath) {
        Files.move(file.absolutePath, filePath)
      }
    }
  }

  files.drop(1).forEach { file ->
    val lastFile = group.last()
    if (ChronoUnit.DAYS.between(lastFile.timeModified, file.timeModified) in 0..3) {
      group.add(file)
    } else {
      moveFilesToDirectory()
      group = mutableListOf(file)
    }
  }

  if (group.size > 0) {
    moveFilesToDirectory()
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
  val goodDirectories = pair.second.first
  val badDirectories = pair.second.second

  uniqueAbsolutePaths.forEach { path ->
    if (path.exists()) {
      categorize(path.toFile(), files, badDirectories, addDirectory)
    }
  }

  deleteDuplicateFiles(files, destinationDirectory).getOrThrow()
  moveFilesAndFilterGoodDirectories(files, goodDirectories, destinationDirectory)

  // There is no need to check if the directory exists before attempting removal.
  badDirectories.asReversed().forEach { directory ->
    directory.delete()
  }
  goodDirectories.forEach { directory ->
    directory.delete()
  }
}
