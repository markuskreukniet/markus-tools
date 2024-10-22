package utils

import org.example.utils.*
import java.nio.file.Files
import java.nio.file.Path
import java.nio.file.Paths
import java.nio.file.attribute.FileTime
import java.time.Instant

fun createFileAndFileSystemFile(directoryPath: String, inputLine: String): Result<FileData> = runCatching {
  val fields = inputLine.split(",")
  val joinedDirectoryPath = Paths.get(directoryPath, fields[0])
  val content = fields[3]
  val name = fields[2]
  val filePath = joinedDirectoryPath.resolve(name)
  val isDirectory = name == ""
  val timeModified = if (fields[1] != "") FileTime.from(Instant.parse(fields[1])) else null

  FileData(
    content = content,
    completeFileInfo = CompleteFileInfo(
      file = filePath.toFile(), // The file is now unusable since the file path is not complete.
      name = name,
      absoluteDirectoryPath = joinedDirectoryPath,
      absolutePath = filePath,
      timeModified = timeModified,
      size = 0L,
      isDirectory = isDirectory,
    )
  )
}

fun createSortedFileSystemFiles(
  rawDelimitedSemicolonString: String
): Result<MutableList<FileData>> = runCatching {
  createSortedFileSystemFiles("", rawDelimitedSemicolonString).getOrThrow()
}

fun createSortedFileSystemFiles(
  directoryPath: String, rawDelimitedSemicolonString: String
): Result<MutableList<FileData>> = runCatching {
  val files = mutableListOf<FileData>()
  val inputLine = mutableListOf<Char>()
  var isCreatingInputLine = false
  val trimmedRawString = rawDelimitedSemicolonString.trim()

  trimmedRawString.forEach { char ->
    if (isCreatingInputLine) {
      if (char != ';') {
        inputLine.add(char)
      } else {
        val file = createFileAndFileSystemFile(directoryPath, inputLine.joinToString("")).getOrThrow()
        files.add(file)
        inputLine.clear()
        isCreatingInputLine = false
      }
    } else if (!char.isWhitespace()) {
      inputLine.add(char)
      isCreatingInputLine = true
    }
  }

  files
}

fun createTemporaryDirectory(): Result<Path> = runCatching {
  Files.createTempDirectory("markus-tools kotlin test_")  // The prefix is optional
}

// Returns the top directory path or a null when it receives only a file name, such as jpg 0.jpg.
fun getTopDirectoryPath(directoryPath: Path): Result<Path?> = runCatching {
  if (directoryPath.nameCount > 0) directoryPath.getName(0) else null
}

fun writeFileAndAddPath(file: FileData, paths: MutableList<Path>): Result<Unit> = runCatching {
  if (!file.completeFileInfo.isDirectory) {
    file.completeFileInfo.absolutePath.toFile().writeText(file.content)
    if (file.completeFileInfo.timeModified != null) {
      Files.setLastModifiedTime(file.completeFileInfo.absolutePath, file.completeFileInfo.timeModified)
    }
  }

  paths.add(file.completeFileInfo.absolutePath)
}

fun writeFilesByMultipleInputs(
  input: String
): Result<Pair<MutableList<Path>?, MutableList<Path>?>> = runCatching {
  if (input.isBlank()) {
    return@runCatching Pair(null, null)
  }

  val files = createSortedFileSystemFiles(input).getOrThrow()

  if (files.size == 0) {
    return@runCatching Pair(null, null)
  }

  val groups = mutableListOf<MutableList<FileData>>(mutableListOf<FileData>(files.first()))
  var previousTopDirectoryPath = getTopDirectoryPath(
    files.first().completeFileInfo.absoluteDirectoryPath
  ).getOrThrow()
  var index = 0

  files.drop(1).forEach { file ->
    val currentTopDirectoryPath = getTopDirectoryPath(file.completeFileInfo.absoluteDirectoryPath).getOrThrow()
    // We can use '==' or '!=' for string-based comparison of the paths.
    if (currentTopDirectoryPath == null || previousTopDirectoryPath != currentTopDirectoryPath) {
      groups.add(mutableListOf<FileData>(file))
      previousTopDirectoryPath = currentTopDirectoryPath
      index++
    } else {
      groups[index].add(file)
    }
  }

  val temporaryDirectories = mutableListOf<Path>()
  val inputPaths = mutableListOf<Path>()

  groups.forEach { group ->
    val directoryPath = createTemporaryDirectory().getOrThrow()
    temporaryDirectories.add(directoryPath)
    var previousDirectoryPath = group.first().completeFileInfo.absoluteDirectoryPath
    group.forEach { file ->
      file.completeFileInfo.absoluteDirectoryPath = directoryPath.resolve(
        file.completeFileInfo.absoluteDirectoryPath
      )
      file.completeFileInfo.absolutePath = directoryPath.resolve(file.completeFileInfo.absolutePath)
      if (file.completeFileInfo.absoluteDirectoryPath != previousDirectoryPath) {
        Files.createDirectories(file.completeFileInfo.absoluteDirectoryPath)
        previousDirectoryPath = file.completeFileInfo.absoluteDirectoryPath
      }
      writeFileAndAddPath(file, inputPaths)
    }
  }

  Pair(temporaryDirectories, inputPaths)
}
