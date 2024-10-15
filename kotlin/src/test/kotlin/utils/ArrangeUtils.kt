package utils

import org.example.utils.*
import java.nio.file.Files
import java.nio.file.Path
import java.nio.file.Paths
import java.time.Instant
import java.time.format.DateTimeFormatter

fun createFileAndFileSystemFile(directoryPath: String, inputLine: String): Result<FileSystemFile> = runCatching {
  val fields = inputLine.split(",")
  val joinedDirectoryPath = Paths.get(directoryPath, fields[0])
  val fileData = fields[3]
  val name = fields[2]
  val filePath = joinedDirectoryPath.resolve(name)
  val isDirectory = name == ""
  val result = if (fields[1] != "") Instant.from(DateTimeFormatter.ISO_DATE_TIME.parse(fields[1])).toEpochMilli()
    else 0L

  FileSystemFile(
    fileData, CompleteFileMetadata(
      name, joinedDirectoryPath.toString(), filePath.toString(), result, 0L, isDirectory, ""
    )
  )
}

fun createSortedFileSystemFiles(rawDelimitedSemicolonString: String): MutableList<FileSystemFile> {
  return createSortedFileSystemFiles("", rawDelimitedSemicolonString)
}

fun createSortedFileSystemFiles(
  directoryPath: String, rawDelimitedSemicolonString: String
): MutableList<FileSystemFile> {
  val files = mutableListOf<FileSystemFile>()
  val inputLine = mutableListOf<Char>()
  var isCreatingInputLine = false
  val trimmedRawString = rawDelimitedSemicolonString.trim()

  trimmedRawString.forEach { char ->
    if (isCreatingInputLine) {
      if (char != ';') {
        inputLine.add(char)
      } else {
        val file = createFileAndFileSystemFile(directoryPath, inputLine.joinToString("")).getOrThrow() // TODO: to string is not efficient
        files.add(file)
        inputLine.clear()
        isCreatingInputLine = false
      }
    } else if (!char.isWhitespace()) {
      inputLine.add(char)
      isCreatingInputLine = true
    }
  }

  return files
}

// TODO: runCatching
fun createTemporaryDirectory(): Path {
  return Files.createTempDirectory("markus-tools kotlin test_")  // The prefix is optional
}

// TODO: runCatching
// Returns the top directory path or a null when it receives only a file name, such as jpg 0.jpg.
fun getTopDirectoryPath(directoryPath: String): Path? {
  val path = Paths.get(directoryPath)
  return if (path.nameCount > 0) path.getName(0) else null
}

fun writeFilesByMultipleInputs(input: String): Pair<MutableList<String>?, MutableList<FileSystemNode>?> {
  if (input.isBlank()) {
    return Pair(null, null)
  }

  val files = createSortedFileSystemFiles(input)

  if (files.size == 0) {
    return Pair(null, null)
  }

  val groups = mutableListOf<MutableList<FileSystemFile>>(mutableListOf<FileSystemFile>(files[0]))
  var previousTopDirectoryPath = getTopDirectoryPath(files[0].completeFileMetadata.absoluteDirectoryPath)
  var index = 0

  files.drop(1).forEach { file ->
    val currentTopDirectoryPath = getTopDirectoryPath(file.completeFileMetadata.absoluteDirectoryPath)
    // We can use '==' or '!=' for string-based comparison.
    if (currentTopDirectoryPath == null || previousTopDirectoryPath != currentTopDirectoryPath) {
      groups.add(mutableListOf<FileSystemFile>(file))
      previousTopDirectoryPath = currentTopDirectoryPath
      index++
    } else {
      groups[index].add(file)
    }
  }

  groups.forEach { group ->
    val directoryPath = createTemporaryDirectory()
  }

  return Pair(mutableListOf<String>(), mutableListOf<FileSystemNode>())
}
