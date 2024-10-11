package org.example.utils

import java.io.File
import java.net.URLConnection

interface FileMetadata {
  val absolutePath: String
  val size: Long
}

data class CompleteFileMetadata(
  val name: String,
  val absoluteDirectoryPath: String,
  override val absolutePath: String,
  val timeModified: Long,
  override val size: Long,
  val isDirectory: Boolean,
  var hash: String
) : FileMetadata

data class FileSystemNode(
  val absolutePath: String,
  val isDirectory: Boolean
)

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

fun filterAndHandleFileMetadata(
  file: File, mode: FileFilterMode, type: FileType, absoluteFilePath: String, handler: (CompleteFileMetadata) -> Unit
): Result<Unit> = runCatching {
  val size: Long = if (file.isFile) file.length() else 0L

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
    absoluteDirectoryPath = "", // TODO:
    absolutePath = absoluteFilePath,
    timeModified = file.lastModified(),
    size = file.length(),
    isDirectory = file.isDirectory,
    hash = ""
  ))
}

fun walkFilterAndHandleFileMetadata(
  absoluteFilePath: String,
  mode: FileFilterMode,
  type: FileType,
  handler: (CompleteFileMetadata) -> Unit
): Result<Unit> = runCatching {
  val rootFile = createExistingFile(absoluteFilePath).getOrThrow() ?: return@runCatching

  if (!rootFile.isFile && !rootFile.isDirectory) {
    return@runCatching
  }

  val files = if (rootFile.isDirectory) rootFile.walk() else sequenceOf(rootFile)
  files.forEach { file ->
    filterAndHandleFileMetadata(file, mode, type, absoluteFilePath, handler).onFailure { throw it }
  }
}

fun createExistingFile(filePath: String): Result<File?> = runCatching {
    val file = File(filePath)
    if (file.exists()) file else null
}
