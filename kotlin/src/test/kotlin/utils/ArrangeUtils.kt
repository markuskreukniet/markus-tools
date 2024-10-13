package utils

import org.example.utils.*
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

fun writeFilesByMultipleInputs(input: String): Pair<MutableList<String>?, MutableList<FileSystemNode>?> {
  if (input.isBlank()) {
    return Pair(null, null)
  }

  return Pair(mutableListOf<String>(), mutableListOf<FileSystemNode>())
}
