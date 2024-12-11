package utils

import org.example.utils.*
import java.nio.file.Files
import java.nio.file.Path
import java.nio.file.Paths
import java.nio.file.attribute.FileTime
import java.time.Instant

fun createFileData(directoryPath: Path?, inputLine: String): Result<FileData> = runCatching {
  val fields = inputLine.split(",")
  val joinedDirectoryPath = Paths.get(directoryPath?.toString() ?: "", fields[0])
  val content = fields[3]
  val name = fields[2]
  val filePath = joinedDirectoryPath.resolve(name)
  val timeModified = if (fields[1] != "") FileTime.from(Instant.parse(fields[1])) else null

  FileData(
    content = content,
    completeFileInfo = CompleteFileInfo(
      file = filePath.toFile(), // The file is now unusable since the file path is not complete.
      absoluteDirectoryPath = joinedDirectoryPath,
      absolutePath = filePath,
      timeModified = timeModified,
      size = 0L // TODO: convert content to size?
    )
  )
}

fun createFilesData(
  rawDelimitedSemicolonString: String
): Result<MutableList<FileData>> = runCatching {
  createFilesData(null, rawDelimitedSemicolonString).getOrThrow()
}

fun createFilesData(
  directoryPath: Path?, rawDelimitedSemicolonString: String
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
        files.add(createFileData(directoryPath, inputLine.joinToString("")).getOrThrow())
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

// Returns the first path segment or a null when it receives only a file name, such as jpg 0.jpg.
fun getFirstPathSegment(directoryPath: Path): Result<Path?> = runCatching {
  if (directoryPath.nameCount > 0) directoryPath.getName(0) else null
}

fun writeFile(file: FileData): Result<Unit> = runCatching {
  if (file.completeFileInfo.file.isDirectory) {
    return@runCatching
  }

  file.completeFileInfo.absolutePath.toFile().writeText(file.content)
  file.completeFileInfo.timeModified?.let { Files.setLastModifiedTime(file.completeFileInfo.absolutePath, it) }
}

fun writeFilesByMultipleInputs(input: String): Result<Pair<MutableList<Path>?, MutableList<Path>?>> = runCatching {
  if (input.isBlank()) {
    return@runCatching Pair(null, null)
  }

  val files = createFilesData(input).getOrThrow()

  if (files.isEmpty()) {
    return@runCatching Pair(null, null)
  }

  files.sortBy { it.completeFileInfo.absolutePath }

  val groups = mutableListOf(mutableListOf(files.first()))
  var previousSegment = getFirstPathSegment(
    files.first().completeFileInfo.absoluteDirectoryPath
  ).getOrThrow()
  var index = 0

  files.drop(1).forEach { file ->
    val currentSegment = getFirstPathSegment(file.completeFileInfo.absoluteDirectoryPath).getOrThrow()
    // We can use '==' or '!=' for string-based comparison of the paths.
    if (currentSegment == null || previousSegment != currentSegment) {
      groups.add(mutableListOf(file))
      previousSegment = currentSegment
      index++
    } else {
      groups[index].add(file)
    }
  }

  // previousSegment = is unnecessary since the possible coming assignments are temporary directory paths,
  // which they were not before.
  val temporaryDirectories = mutableListOf<Path>()
  val inputPaths = mutableListOf<Path>()

  groups.forEach { group ->
    val directoryPath = createTemporaryDirectory().getOrThrow()
    temporaryDirectories.add(directoryPath)
    group.forEach { file ->
      file.completeFileInfo.absoluteDirectoryPath = directoryPath.resolve(
        file.completeFileInfo.absoluteDirectoryPath
      )
      file.completeFileInfo.absolutePath = directoryPath.resolve(file.completeFileInfo.absolutePath)
      if (file.completeFileInfo.absoluteDirectoryPath != previousSegment) {
        Files.createDirectories(file.completeFileInfo.absoluteDirectoryPath)
        previousSegment = file.completeFileInfo.absoluteDirectoryPath
      }
      writeFile(file)
      inputPaths.add(file.completeFileInfo.absolutePath)
    }
  }

  Pair(temporaryDirectories, inputPaths)
}
