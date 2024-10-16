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
  val timeModified = if (fields[1] != "") Instant.from(DateTimeFormatter.ISO_DATE_TIME.parse(fields[1])).toEpochMilli()
    else 0L

  FileSystemFile(
    fileData, CompleteFileMetadata(
      name = name,
      absoluteDirectoryPath = joinedDirectoryPath,
      absolutePath = filePath,
      timeModified = timeModified,
      size = 0L,
      isDirectory = isDirectory,
      hash = ""
    )
  )
}

fun createSortedFileSystemFiles(
  rawDelimitedSemicolonString: String
): Result<MutableList<FileSystemFile>> = runCatching {
  createSortedFileSystemFiles("", rawDelimitedSemicolonString).getOrThrow()
}

fun createSortedFileSystemFiles(
  directoryPath: String, rawDelimitedSemicolonString: String
): Result<MutableList<FileSystemFile>> = runCatching {
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

  files
}

fun createTemporaryDirectory(): Result<Path> = runCatching {
  Files.createTempDirectory("markus-tools kotlin test_")  // The prefix is optional
}

// Returns the top directory path or a null when it receives only a file name, such as jpg 0.jpg.
fun getTopDirectoryPath(directoryPath: Path): Result<Path?> = runCatching {
  if (directoryPath.nameCount > 0) directoryPath.getName(0) else null
}

//func testingIfFileWriteItAndAppendFileSystemNode(t *testing.T, file FileSystemFile, nodes *[]FileSystemNode) {
//  t.Helper()
//
//  if !file.FileMetadata.IsDirectory {
//    TestingWriteFile(t, file.FileMetadata.Path, file.Data)
//    if !file.FileMetadata.TimeModified.IsZero() {
//      if err := os.Chtimes(file.FileMetadata.Path, time.Now(), file.FileMetadata.TimeModified); err != nil {
//      t.Errorf("Failed to change the access and modification times of the file: %v", err)
//    }
//    }
//  }
//
//  *nodes = append(*nodes, FileSystemNode{
//      Path:        file.FileMetadata.Path,
//      IsDirectory: file.FileMetadata.IsDirectory,
//  })
//}

fun asdf(file: FileSystemFile, paths: MutableList<Path>) {
  if (!file.completeFileMetadata.isDirectory) {
    // File(filePath).writeText(content)
    if (file.completeFileMetadata.timeModified > 0L) {

    }
  }

  paths.add(file.completeFileMetadata.absolutePath)
}

fun writeFilesByMultipleInputs(
  input: String
): Result<Pair<MutableList<String>?, MutableList<Path>?>> = runCatching {
  if (input.isBlank()) {
    return@runCatching Pair(null, null)
  }

  val files = createSortedFileSystemFiles(input).getOrThrow()

  if (files.size == 0) {
    return@runCatching Pair(null, null)
  }

  val groups = mutableListOf<MutableList<FileSystemFile>>(mutableListOf<FileSystemFile>(files.first()))
  var previousTopDirectoryPath = getTopDirectoryPath(
    files.first().completeFileMetadata.absoluteDirectoryPath
  ).getOrThrow()
  var index = 0

  files.drop(1).forEach { file ->
    val currentTopDirectoryPath = getTopDirectoryPath(file.completeFileMetadata.absoluteDirectoryPath).getOrThrow()
    // We can use '==' or '!=' for string-based comparison of the paths.
    if (currentTopDirectoryPath == null || previousTopDirectoryPath != currentTopDirectoryPath) {
      groups.add(mutableListOf<FileSystemFile>(file))
      previousTopDirectoryPath = currentTopDirectoryPath
      index++
    } else {
      groups[index].add(file)
    }
  }

  val temporaryDirectories = mutableListOf<Path>()

  groups.forEach { group ->
    val directoryPath = createTemporaryDirectory().getOrThrow()
    temporaryDirectories.add(directoryPath)
    var previousDirectoryPath = group.first().completeFileMetadata.absoluteDirectoryPath
    group.forEach { file ->
      file.completeFileMetadata.absoluteDirectoryPath = directoryPath.resolve(
        file.completeFileMetadata.absoluteDirectoryPath
      )
      file.completeFileMetadata.absolutePath = directoryPath.resolve(file.completeFileMetadata.absolutePath)
      if (file.completeFileMetadata.absoluteDirectoryPath != previousDirectoryPath) {
        Files.createDirectory(file.completeFileMetadata.absoluteDirectoryPath)
        previousDirectoryPath = file.completeFileMetadata.absoluteDirectoryPath
      }
    }
  }

  Pair(mutableListOf<String>(), mutableListOf<Path>())
}
