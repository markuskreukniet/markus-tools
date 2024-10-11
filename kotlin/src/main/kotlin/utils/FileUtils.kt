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
): Result<Unit> {
  val size: Long = if (file.isFile) file.length() else 0L

  // is file check
  if (file.isFile && mode == FileFilterMode.DIRECTORIES) {
    return Result.success(Unit)
  }

  // is directory check
  if (file.isDirectory && (mode == FileFilterMode.FILES || mode == FileFilterMode.NON_ZERO_BYTE_FILES)) {
    return Result.success(Unit)
  }

  // is zero byte file check
  if (file.isFile && size == 0L &&
    (mode == FileFilterMode.NON_ZERO_BYTE_FILES || mode == FileFilterMode.NON_ZERO_BYTE_FILES_AND_DIRECTORIES)) {
    return Result.success(Unit)
  }

  // is text file check
  val isTextFile = isTextFile(file).getOrThrow()
  if (type == FileType.TEXT_FILES && isTextFile) {
    return Result.success(Unit)
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

  return Result.success(Unit)
}

fun walkFilterAndHandleFileMetadata(
  absoluteFilePath: String,
  mode: FileFilterMode,
  type: FileType,
  handler: (CompleteFileMetadata) -> Unit
): Result<Unit> {
  val rootFile = createExistingFile(absoluteFilePath).getOrThrow() ?: return Result.success(Unit)

  if (!rootFile.isFile && !rootFile.isDirectory) {
    return Result.success(Unit)
  }

  val files = if (rootFile.isDirectory) rootFile.walk() else sequenceOf(rootFile)
  files.forEach { file ->
    filterAndHandleFileMetadata(file, mode, type, absoluteFilePath, handler).getOrThrow()
  }

  return Result.success(Unit)
}

fun createExistingFile(filePath: String): Result<File?> = runCatching {
    val file = File(filePath)
    if (file.exists()) file else null
}
