package utils

import org.example.utils.*
import java.nio.file.Files
import java.nio.file.Path
import java.nio.file.Paths
import java.nio.file.attribute.FileTime
import java.time.Instant

fun createFileData(directoryPath: String, inputLine: String): Result<FileData> = runCatching {
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
      absoluteDirectoryPath = joinedDirectoryPath,
      absolutePath = filePath,
      timeModified = timeModified,
      size = 0L // TODO: convert content to size?
    )
  )
}

fun createFileSystemFiles(
  rawDelimitedSemicolonString: String
): Result<MutableList<FileData>> = runCatching {
  createFileSystemFiles("", rawDelimitedSemicolonString).getOrThrow()
}

fun createFileSystemFiles(
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
        val file = createFileData(directoryPath, inputLine.joinToString("")).getOrThrow()
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

fun writeFile(file: FileData): Result<Unit> = runCatching {
  if (file.completeFileInfo.file.isDirectory) {
    return@runCatching
  }

  file.completeFileInfo.absolutePath.toFile().writeText(file.content)
  file.completeFileInfo.timeModified?.let { Files.setLastModifiedTime(file.completeFileInfo.absolutePath, it) }
}

fun resolveAbsolutePaths(directoryPath: Path, file: FileData) {
  file.completeFileInfo.absoluteDirectoryPath = directoryPath.resolve(
    file.completeFileInfo.absoluteDirectoryPath
  )
  file.completeFileInfo.absolutePath = directoryPath.resolve(file.completeFileInfo.absolutePath)
}

fun writeFilesBySingleInput(input: String): Result<Path?> = runCatching {
  if (input.isBlank()) {
    return@runCatching null
  }

  val files = createFileSystemFiles(input).getOrThrow()

  if (files.isEmpty()) {
    return@runCatching null
  }

  val directoryPath = createTemporaryDirectory().getOrThrow()

  files.forEach { file ->
    resolveAbsolutePaths(directoryPath, file)
    writeFile(file)
  }

  directoryPath
}

fun writeFilesByMultipleInputs(input: String): Result<Pair<MutableList<Path>?, MutableList<Path>?>> = runCatching {
  if (input.isBlank()) {
    return@runCatching Pair(null, null)
  }

  val files = createFileSystemFiles(input).getOrThrow()

  if (files.isEmpty()) {
    return@runCatching Pair(null, null)
  }

  files.sortBy { it.completeFileInfo.absolutePath }

  val groups = mutableListOf(mutableListOf(files.first()))
  var previousTopDirectoryPath = getTopDirectoryPath(
    files.first().completeFileInfo.absoluteDirectoryPath
  ).getOrThrow()
  var index = 0

  files.drop(1).forEach { file ->
    val currentTopDirectoryPath = getTopDirectoryPath(file.completeFileInfo.absoluteDirectoryPath).getOrThrow()
    // We can use '==' or '!=' for string-based comparison of the paths.
    if (currentTopDirectoryPath == null || previousTopDirectoryPath != currentTopDirectoryPath) {
      groups.add(mutableListOf(file))
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
      resolveAbsolutePaths(directoryPath, file)
      if (file.completeFileInfo.absoluteDirectoryPath != previousDirectoryPath) {
        Files.createDirectories(file.completeFileInfo.absoluteDirectoryPath)
        previousDirectoryPath = file.completeFileInfo.absoluteDirectoryPath
      }
      writeFile(file)
      inputPaths.add(file.completeFileInfo.absolutePath)
    }
  }

  Pair(temporaryDirectories, inputPaths)
}
