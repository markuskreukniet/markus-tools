package org.example.utils

import java.io.File
import java.net.URLConnection
import java.nio.file.Path
import java.nio.file.attribute.FileTime
import kotlin.io.path.getLastModifiedTime

data class FileSystemFile(
  val data: String,
  val completeFileMetadata: CompleteFileMetadata
)

interface FileMetadata {
  val absolutePath: Path
  val size: Long
}

data class CompleteFileMetadata(
  val name: String,
  var absoluteDirectoryPath: Path,
  override var absolutePath: Path,
  val timeModified: FileTime?,
  override val size: Long,
  val isDirectory: Boolean,
  var hash: String
) : FileMetadata

enum class FileFilterMode {
  FILES,
  NON_ZERO_BYTE_FILES,
  FILES_AND_DIRECTORIES,
  NON_ZERO_BYTE_FILES_AND_DIRECTORIES,
  DIRECTORIES
}

enum class FileType {
  ALL_FILES,
  TEXT_FILES
}

fun isTextFile(file: File): Result<Boolean> = runCatching {
  val mimeType = file.inputStream().use { inputStream ->
    URLConnection.guessContentTypeFromStream(inputStream)
  }

  mimeType?.startsWith("text") == true
}

fun getDirectoryPath(filePath: Path, isDirectory: Boolean): Path {
  return if (isDirectory || filePath.parent == null) {
    filePath
  } else {
    filePath.parent
  }
}

fun filterAndHandleFileMetadata(
  file: File, mode: FileFilterMode, type: FileType, absoluteFilePath: Path, handler: (CompleteFileMetadata) -> Unit
): Result<Unit> = runCatching {
  val size = if (file.isFile) file.length() else 0L

  // is file check
  if (file.isFile && mode == FileFilterMode.DIRECTORIES) {
    return@runCatching
  }

  // is directory check
  if (file.isDirectory && (mode == FileFilterMode.FILES || mode == FileFilterMode.NON_ZERO_BYTE_FILES)) {
    return@runCatching
  }

  // is zero byte file check
  if (file.isFile && size == 0L &&
    (mode == FileFilterMode.NON_ZERO_BYTE_FILES || mode == FileFilterMode.NON_ZERO_BYTE_FILES_AND_DIRECTORIES)) {
    return@runCatching
  }

  // is text file check
  val isTextFile = isTextFile(file).getOrThrow()
  if (type == FileType.TEXT_FILES && isTextFile) {
    return@runCatching
  }

  handler(CompleteFileMetadata(
    name = file.name,
    absoluteDirectoryPath = getDirectoryPath(absoluteFilePath, file.isDirectory),
    absolutePath = absoluteFilePath,
    timeModified = absoluteFilePath.getLastModifiedTime(),
    size = file.length(),
    isDirectory = file.isDirectory,
    hash = ""
  ))
}

fun walkFilterAndHandleFileMetadata(
  absoluteFilePath: Path,
  mode: FileFilterMode,
  type: FileType,
  handler: (CompleteFileMetadata) -> Unit
): Result<Unit> = runCatching {
  val rootFile = absoluteFilePath.toFile()

  if (!rootFile.isFile && !rootFile.isDirectory) {
    return@runCatching
  }

  val files = if (rootFile.isDirectory) rootFile.walk() else sequenceOf(rootFile)
  files.forEach { file ->
    filterAndHandleFileMetadata(file, mode, type, absoluteFilePath, handler).onFailure { throw it }
  }
}
