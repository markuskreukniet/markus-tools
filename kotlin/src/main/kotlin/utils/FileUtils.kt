package org.example.utils

import java.io.File
import java.net.URLConnection
import java.nio.file.Path
import java.nio.file.attribute.FileTime
import java.time.Instant
import kotlin.io.path.getLastModifiedTime

data class FileData(
  val content: String,
  val completeFileInfo: CompleteFileInfo
)

interface FileInfo {
  val file: File
  val size: Long
  val absolutePath: Path
}

interface DuplicateFileInfo {
  val file: File
  val size: Long
  val absolutePath: Path
}

data class DFFileInfo(
  override val file: File,
  override val size: Long,
  override val absolutePath: Path
) : DuplicateFileInfo

// TODO: rename to DateRangeFileInfo
data class FTDRFileInfo(
  override val file: File,
  override val size: Long,
  override val absolutePath: Path,
  val timeModified: Instant
) : DuplicateFileInfo

data class MinimalFileInfo(
  override val file: File,
  override val size: Long,
  override val absolutePath: Path
) : FileInfo

data class CompleteFileInfo(
  override val file: File,
  val name: String,
  var absoluteDirectoryPath: Path,
  override var absolutePath: Path,
  val timeModified: FileTime?,
  override val size: Long,
  val isDirectory: Boolean, // TODO: is isDirectory needed?
) : FileInfo

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

fun resolveDirectoryPath(filePath: Path, isDirectory: Boolean): Path {
  return if (isDirectory || filePath.parent == null) {
    filePath
  } else {
    filePath.parent
  }
}

fun filterAndHandleFileInfo(
  file: File, mode: FileFilterMode, type: FileType, absoluteFilePath: Path, handler: (FileInfo) -> Unit
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
  if (type == FileType.TEXT_FILES) {
    if (!isTextFile(file).getOrThrow()) {
      return@runCatching
    }
  }

  handler(CompleteFileInfo(
    file = file,
    name = file.name,
    absoluteDirectoryPath = resolveDirectoryPath(absoluteFilePath, file.isDirectory),
    absolutePath = absoluteFilePath,
    timeModified = absoluteFilePath.getLastModifiedTime(),
    size = file.length(),
    isDirectory = file.isDirectory,
  ))
}

fun walkFilterAndHandleFileInfo(
  absoluteFilePath: Path,
  mode: FileFilterMode,
  type: FileType,
  handler: (FileInfo) -> Unit
): Result<Unit> = runCatching {
  val rootFile = absoluteFilePath.toFile()

  if (!rootFile.isFile && !rootFile.isDirectory) {
    return@runCatching
  }

  val files = if (rootFile.isDirectory) rootFile.walk() else sequenceOf(rootFile)
  files.forEach { file ->
    filterAndHandleFileInfo(file, mode, type, absoluteFilePath, handler).onFailure { throw it }
  }
}
