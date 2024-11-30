package org.example

import org.example.utils.FDateRangeFileInfo
import org.example.utils.createDuplicateFileInfoGroupsByHash
import java.io.File
import java.nio.file.Files
import java.nio.file.Path
import java.nio.file.Paths
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

  val parseDate = fun(rawDate: String): Result<LocalDate> = runCatching {
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
): Result<Pair<MutableList<FDateRangeFileInfo>, Pair<MutableSet<File>, MutableList<File>>>> = runCatching {
  val files = mutableListOf<FDateRangeFileInfo>()
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

  val directories: Collection<File> = goodDirectories + badDirectories

  directories.forEach { directory ->
    directory.walk().drop(1).forEach { file ->
      categorize(file, files, badDirectories, addDirectory)
    }
  }

  Pair(files, Pair(goodDirectories, badDirectories))
}

fun categorize(
  file: File,
  files: MutableList<FDateRangeFileInfo>,
  badDirectories: MutableList<File>,
  handler: (MutableList<File>, File) -> Unit
) = runCatching {
  if (file.isDirectory) {
    handler(badDirectories, file)
  } else if (file.isFile) {
    val size = file.length()
    if (size > 0L) {
      val path = file.toPath()
      files.add(FDateRangeFileInfo(
        file = file,
        size = size,
        path = path,
        timeModified = path.getLastModifiedTime().toInstant(),
        newName = null
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
): Result<List<(MutableList<FDateRangeFileInfo>, MutableList<File>) ->
Result<MutableList<FDateRangeFileInfo>>>> = runCatching {
  val addAllFilesInfo = fun(files: MutableList<File>, filesInfo: List<FDateRangeFileInfo>) {
    filesInfo.forEach { file ->
      files.add(file.file)
    }
  }

  val addBadFilesInfoAndReplaceGoodFiles = fun(
    badFiles: MutableList<File>, goodFiles: MutableList<FDateRangeFileInfo>, file: FDateRangeFileInfo
  ) {
    addAllFilesInfo(badFiles, goodFiles)
    goodFiles.clear()
    goodFiles.add(file)
  }

  val categorizeOnShortestFileNameLength = fun(
    files: MutableList<FDateRangeFileInfo>, badFiles: MutableList<File>
  ): Result<MutableList<FDateRangeFileInfo>> = runCatching {
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

    good
  }

  val categorizeOnValidDateRangeDirectoryName = fun(
    files: MutableList<FDateRangeFileInfo>, badFiles: MutableList<File>
  ): Result<MutableList<FDateRangeFileInfo>> {
    val tempGood1Files = mutableListOf<FDateRangeFileInfo>()
    val tempGood2Files = mutableListOf<FDateRangeFileInfo>()
    val tempBadFiles = mutableListOf<FDateRangeFileInfo>()

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
      return Result.success(tempGood2Files)
    }

    if (tempGood1Files.isNotEmpty()) {
      addAllFilesInfo(badFiles, tempBadFiles)
      return Result.success(tempGood1Files)
    }

    return Result.success(tempBadFiles)
  }

  val categorizeOnNewestTimeModified = fun(
    files: MutableList<FDateRangeFileInfo>, badFiles: MutableList<File>
  ): Result<MutableList<FDateRangeFileInfo>> = runCatching {
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

    good
  }

  val categorizeOnFirstFile = fun(
    files: MutableList<FDateRangeFileInfo>, badFiles: MutableList<File>
  ): Result<MutableList<FDateRangeFileInfo>> = runCatching {
    val good = mutableListOf(files.first())

    addAllFilesInfo(badFiles, files.drop(1))

    good
  }

  listOf(
    categorizeOnShortestFileNameLength,
    categorizeOnValidDateRangeDirectoryName,
    categorizeOnNewestTimeModified,
    categorizeOnFirstFile
  )
}

fun deleteDuplicateFiles(
  files: MutableList<FDateRangeFileInfo>, destinationDirectory: File
): Result<Unit> = runCatching {
  val groups = createDuplicateFileInfoGroupsByHash(files, false).getOrThrow() ?: return@runCatching
  val handlers = createHandlers(destinationDirectory).getOrThrow()
  val badFiles = mutableListOf<File>()

  files.clear()

  groups.forEachIndexed { index, group ->
    for (handler in handlers) {
      // group and groups[index] are different references
      if (groups[index].size > 1) {
        groups[index] = handler(group, badFiles).getOrThrow()
      } else {
        files.add(groups[index].first())
        break
      }
    }
  }

  badFiles.forEach { file ->
    file.delete()
  }
}

fun moveFilesAndFilterGoodDirectories(
  files: MutableList<FDateRangeFileInfo>, goodDirectories: MutableSet<File>, destinationDirectory: File
) = runCatching {
  if (files.isEmpty()) {
    return@runCatching
  }

  files.sortBy { it.timeModified }

  // In Kotlin, a mutableSet provides O(1) access time and is ordered.
  // However, we still need to be able to look up file names, so we use a mutableSet and a mutableList.
  var fileNames = mutableSetOf<String>()
  var group = mutableListOf<FDateRangeFileInfo>()

  val replaceFileNamesAndGroup = fun(file: FDateRangeFileInfo) {
    fileNames = mutableSetOf(file.file.name)
    group = mutableListOf(file)
  }

  replaceFileNamesAndGroup(files.first())

  val formatTimeModified = fun(file: FDateRangeFileInfo): Result<String> = runCatching {
    file.timeModified.atZone(ZoneId.systemDefault())
      .toLocalDate()
      .format(getDateTimeFormatter().getOrThrow())
  }

  val moveFilesToDirectory = fun() = runCatching {
    val firstFile = group.first()
    val lastFile = group.last()

    var directoryName = formatTimeModified(firstFile).getOrThrow()
    if (ChronoUnit.DAYS.between(firstFile.timeModified, lastFile.timeModified) >= 1) {
      directoryName += " - ${formatTimeModified(lastFile).getOrThrow()}"
    }

    val joinedDirectory = File(destinationDirectory.path, directoryName)
    if (joinedDirectory in goodDirectories) {
      goodDirectories.remove(joinedDirectory)
    } else {
      joinedDirectory.mkdir()
    }

    group.forEach { file ->
      if (file.newName == null) {
        val filePath = Paths.get(joinedDirectory.toString(), file.file.name)
        if (filePath != file.path) {
          Files.move(file.path, filePath)
        }
      } else {
        Files.move(file.path, Paths.get(joinedDirectory.toString(), file.newName))
      }
    }
  }

  files.drop(1).forEach { file ->
    val lastFile = group.last()
    if (ChronoUnit.DAYS.between(lastFile.timeModified, file.timeModified) in 0..3) {
      if (file.file.name in fileNames) {
        var disambiguationNumber = 2
        while (disambiguationNumber <= 9) {
          val name = "$file.file.nameWithoutExtension ${disambiguationNumber}.$file.file.extension"
          file.newName = name
          if (name !in fileNames) {
            fileNames.add(name)
            break
          }
          disambiguationNumber++
        }
        if (disambiguationNumber == 9) {
          // TODO: exception
        }
      } else {
        fileNames.add(file.file.name)
      }
      group.add(file)
    } else {
      moveFilesToDirectory()
      replaceFileNamesAndGroup(file)
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

  val pair = categorizeFilesAndDirectories(destinationDirectory).getOrThrow()
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
