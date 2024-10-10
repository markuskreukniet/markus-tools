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

fun isTextFile(file: File): Boolean {
  val mimeType = file.inputStream().use { inputStream ->
    URLConnection.guessContentTypeFromStream(inputStream)
  }

  return mimeType?.startsWith("text") == true
}

fun filterAndHandleFileMetadata(
  file: File, mode: FileFilterMode, type: FileType, absoluteFilePath: String, handler: (CompleteFileMetadata) -> Unit) {
  val size: Long = if (file.isFile) file.length() else 0L

  // is file check
  if (file.isFile && mode == FileFilterMode.DIRECTORIES) {
    return
  }

  // is directory check
  if (file.isDirectory && (mode == FileFilterMode.FILES || mode == FileFilterMode.NON_ZERO_BYTE_FILES)) {
    return
  }

  // is zero byte file check
  if (file.isFile && size == 0L &&
    (mode == FileFilterMode.NON_ZERO_BYTE_FILES || mode == FileFilterMode.NON_ZERO_BYTE_FILES_AND_DIRECTORIES)) {
    return
  }

  // is text file check
  if (type == FileType.TEXT_FILES && isTextFile(file)) {
    return
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
): Result<Unit> {
  val rootFile = createExistingFile(absoluteFilePath)
    .getOrElse { return Result.failure(it) } ?: return Result.success(Unit)

  if (rootFile.isFile) {
    filterAndHandleFileMetadata(rootFile, mode, type, absoluteFilePath, handler)
  } else if (rootFile.isDirectory) {
    rootFile.walk().forEach { file ->
      filterAndHandleFileMetadata(file, mode, type, absoluteFilePath, handler)
    }
  }

  return Result.success(Unit)
}

fun createExistingFile(filePath: String): Result<File?> {
  return runCatching {
    val file = File(filePath)
    if (file.exists()) file else null
  }
}
